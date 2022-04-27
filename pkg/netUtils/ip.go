package netUtils

import (
	"github.com/kostua16/go_simple_logger/pkg/logger"
	"net"
)

var log = logger.CreateLogger("net.utils")

const (
	WSLNetwork        = "vEthernet (WSL)"
	VirtualBoxNetwork = "VirtualBox Host-Only Network"
)

func isWSLNetwork(network net.Interface) bool {
	return network.Name == WSLNetwork
}

func isVirtualBoxNetwork(network net.Interface) bool {
	return network.Name == VirtualBoxNetwork
}

// GetLocalIP Get first local ip of this machine
func GetLocalIP() net.IP {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			if i.Flags&net.FlagLoopback != 0 { // localhost
				log.Debugf("GetLocalIP: found iface: %s, %s ", i.Name, i.Flags)
				addrs, err := i.Addrs()
				if err == nil {
					for _, addr := range addrs {
						if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
							log.Debugf("GetLocalIP: found ip: %s", ipnet.IP.String())
							return ipnet.IP.To4()
						}
					}
				}
			}

		}
	}
	return net.IPv4(127, 0, 0, 1)
}

// GetOutboundIP Get first outbound ip of this machine
func GetOutboundIP() net.IP {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			log.Debugf("GetOutboundIP: found iface: %s, %s ", i.Name, i.Flags)
			if i.Flags&net.FlagUp == 0 { // is down
				continue
			}
			if i.Flags&net.FlagLoopback != 0 { // localhost
				continue
			}
			if isWSLNetwork(i) || isVirtualBoxNetwork(i) { // Virtual networks
				continue
			}
			addrs, err := i.Addrs()
			if err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
						log.Debugf("GetOutboundIP: found ip: %s", ipnet.IP.String())
						return ipnet.IP.To4()
					}
				}
			}

		}
	}
	return GetLocalIP()
}

// GetExternalIP Get preferred outbound ip of this machine
func GetExternalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return GetOutboundIP()
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
