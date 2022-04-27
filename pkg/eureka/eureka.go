package eureka

import (
	"errors"
	"github.com/hudl/fargo"
	"github.com/kostua16/go_simple_logger/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

var log = logger.CreateLogger("eureka")

type InstanceInfo struct {
	config             *fargo.Config
	instance           *fargo.Instance
	connection         *fargo.EurekaConnection
	stopInstance       chan bool
	instanceRegistered bool
	serverOK           bool
}

func (i *InstanceInfo) ConfigureInstance(host string, port int, ip string, name string) {
	log.Infof("Creating eureka instance for %s:%s (%s) = %s", host, port, ip, name)
	instance := fargo.Instance{
		HostName:         host,
		Port:             port,
		App:              name,
		IPAddr:           ip,
		VipAddress:       ip,
		SecureVipAddress: ip,
		Status:           fargo.UP,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		PortEnabled:      true,
		LeaseInfo: fargo.LeaseInfo{
			DurationInSecs: 90,
		},
		InstanceId: name + "-" + host + "-" + strconv.Itoa(port),
	}
	log.Infof("Created eureka instance for %s:%s (%s) = %s", host, port, ip, name)
	i.stopInstance = make(chan bool)
	i.instance = &instance
}

func (i *InstanceInfo) StartInstance() bool {
	if i.IsConnected() && i.InstanceConfigured() {
		i.stopInstance = make(chan bool)
		i.sentInstanceHeartbeat()
		go i.sentInstanceHeartbeats()
		return true
	}
	return false
}

func (i *InstanceInfo) StopInstance() bool {
	if i.IsConnected() && i.InstanceConfigured() && i.InstanceRegistered() {
		i.deregisterInstance()
		return true
	}
	return false
}

func (i *InstanceInfo) registerCurrentInstance() {
	if i.IsConnected() && i.InstanceConfigured() {
		log.Infof("Registering instance %s..", i.instance.Id())
		regErr := i.connection.ReregisterInstance(i.instance)
		if regErr != nil {
			log.Errorf("Failed to register instance %s: %s", i.instance.Id(), regErr)
			i.instanceRegistered = false
		} else {
			i.instanceRegistered = true
			log.Infof("Registered instance %s.", i.instance.Id())
		}
	}
}

func (i *InstanceInfo) deregisterInstance() {
	if i.IsConnected() && i.InstanceConfigured() && i.InstanceRegistered() {
		log.Infof("De-registering instance %s..", i.instance.Id())
		deRegErr := i.connection.DeregisterInstance(i.instance)
		if deRegErr != nil {
			log.Errorf("Failed to de-register instance %s: %s", i.instance.Id(), deRegErr)
		} else {
			log.Infof("De-registered instance %s.", i.instance.Id())
		}
		i.stopInstance <- true
		i.instanceRegistered = false
	}
}

func (i *InstanceInfo) ChangeInstanceStatus(status fargo.StatusType) error {
	if i.IsConnected() && i.InstanceConfigured() && i.InstanceRegistered() {
		log.Infof("Updating instance %s status to %s", i.instance.Id(), status)
		updErr := i.connection.UpdateInstanceStatus(i.instance, status)
		if updErr != nil {
			log.Errorf("Failed to update instance %s status: %s", i.instance.Id(), updErr)
		} else {
			log.Infof("Updated instance %s status to %s", i.instance.Id(), status)
		}
		return updErr
	} else {
		return errors.New("instance is not connected")
	}
}

func (i *InstanceInfo) sentInstanceHeartbeat() {
	if i.InstanceConfigured() && i.IsConnected() {
		log.Infof("Sending heartbeat for instance %s...", i.instance.Id())
		hbErr := i.connection.HeartBeatInstance(i.instance)
		if hbErr != nil {
			code, ok := fargo.HTTPResponseStatusCode(hbErr)
			if ok && code == http.StatusNotFound {
				log.Infof("Heartbeat shown that instance %s were not found, going to register it", i.instance.Id())
				i.registerCurrentInstance()
			} else {
				log.Errorf("Failed to send heartbeat for instance %s: %s", i.instance.Id(), hbErr)
			}
		} else {
			log.Infof("Sent heartbeat for instance %s.", i.instance.Id())
		}
	}
}

func (i *InstanceInfo) sentInstanceHeartbeats() {
	i.registerCurrentInstance()
	ticker := time.NewTicker(time.Duration(i.config.Eureka.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	i.stopInstance = make(chan bool)
	for {
		select {
		case <-i.stopInstance:
			return
		case <-ticker.C:
			i.sentInstanceHeartbeat()
		}
	}
}

func (i *InstanceInfo) initConfig() {
	cfg := fargo.Config{}
	cfg.Eureka.ServiceUrls = []string{}
	cfg.Eureka.PollIntervalSeconds = 30
	cfg.Eureka.Retries = 3
	cfg.Eureka.ConnectTimeoutSeconds = 10
	i.config = &cfg
}

func (i *InstanceInfo) SetDefaultServerUrl() {
	i.SetServerUrl("http://localhost:8761/eureka/eureka")
}

func (i *InstanceInfo) SetServerUrl(eurekaUrl string) {
	log.Infof("Creating eureka client for %s...", eurekaUrl)
	i.config.Eureka.ServiceUrls = []string{eurekaUrl}
	eureka := fargo.NewConnFromConfig(*i.config)
	log.Infof("Created eureka client for %s.", eurekaUrl)
	i.connection = &eureka
	i.pingServer()
	go i.sentPings()
}

func (i *InstanceInfo) pingServer() bool {
	if i.IsClientConfigured() {
		url := i.connection.SelectServiceURL()
		log.Debugf("Sending ping to eureka server %s..", url)
		response, err := http.Get(url)
		if err != nil {
			log.Debugf("Ping to eureka server %s = FAILED (err: %s)", url, err.Error())
			return false
		}
		if response.StatusCode != http.StatusOK {
			log.Debugf("Ping to eureka server %s = FAILED (code: %s)", url, strconv.Itoa(response.StatusCode))
			return false
		}
		log.Debugf("Ping to eureka server %s = OK", url)
		return true
	}
	return false
}

func (i *InstanceInfo) sentPings() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	i.stopInstance = make(chan bool)
	for {
		select {
		case <-ticker.C:
			if i.IsClientConfigured() {
				i.serverOK = i.pingServer()
			} else {
				return
			}

		}
	}
}

func (i *InstanceInfo) IsClientConfigured() bool {
	return i.connection != nil && len(i.config.Eureka.ServiceUrls) > 0
}

func (i *InstanceInfo) IsConnected() bool {
	return i.IsClientConfigured() && i.serverOK
}

func (i *InstanceInfo) InstanceConfigured() bool {
	return i.instance != nil
}

func (i *InstanceInfo) InstanceRegistered() bool {
	return i.instanceRegistered
}

func (i *InstanceInfo) AppId() string {
	if i.instance != nil {
		return i.instance.Id()
	}
	return ""
}

func NewEmptyInstance() *InstanceInfo {
	var data = &InstanceInfo{}
	data.initConfig()
	return data
}

func NewClientOnlyInstance(eurekaUrl string) *InstanceInfo {
	var data = NewEmptyInstance()
	data.SetServerUrl(eurekaUrl)
	return data
}

func NewInstance(eurekaUrl string, host string, port int, ip string, name string) *InstanceInfo {
	var data = NewClientOnlyInstance(eurekaUrl)
	data.ConfigureInstance(host, port, ip, name)
	return data
}
