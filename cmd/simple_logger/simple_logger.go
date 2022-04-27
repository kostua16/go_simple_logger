package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kostua16/go_simple_logger/pkg/api/logApi/logApi.v1"
	"github.com/kostua16/go_simple_logger/pkg/db"
	"github.com/kostua16/go_simple_logger/pkg/eureka"
	"github.com/kostua16/go_simple_logger/pkg/logger"
	"github.com/kostua16/go_simple_logger/pkg/netUtils"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	log       = logger.CreateLogger("main")
	eurekaUrl = flag.String("eureka", "http://127.0.0.1:8761", "Eureka server")
	hostname  = flag.String("host", netUtils.GetHostName(), "Public hostname for this environment")
	port      = flag.Int("port", 8080, "Public port to start listening")
	ip        = flag.String("ip", netUtils.GetExternalIP().String(), "Eureka server")
	appName   = flag.String("name", "logger", "The name of the app in the Eureka")
	help      = flag.Bool("help", false, "Show help")
	verbose   = flag.String("verbose", "WARNING", "Logging level: CRITICAL,ERROR,WARNING,NOTICE,INFO,DEBUG")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	logger.SetLevelStr(*verbose)
	logger.SetModuleLevelStr("CRITICAL", "fargo")
	fmt.Println("Starting hello-world server...")

	eurekaInst := eureka.NewInstance(*eurekaUrl, *hostname, *port, *ip, *appName)

	if eurekaInst.StartInstance() {
		defer eurekaInst.StopInstance()
	}

	dbConn := db.NewConnection("simple_logger.db")
	dbOpenErr := dbConn.Open()
	if dbOpenErr != nil {
		log.Panicf("Failed to create database connection: %v", dbOpenErr)
	} else {
		defer func(dbConn *db.Connection) {
			err := dbConn.Close()
			if err != nil {
				log.Panicf("Failed to close database connection: %v", err)
			}
		}(dbConn)
	}

	logApi, err := logApi_v1.CreateApi(dbConn)

	if err != nil {
		log.Panicf("Failed to create logApi: %v", err)
	}

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/api/logs/v1", logApi)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})

	//http.HandleFunc("/", helloServer)

	fmt.Println("Started hello-world server.")

	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		panic(err)
	}
	fmt.Println("Finished hello-world server.")
}
