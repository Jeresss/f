package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
)

const serviceName = "Medium service"
const serviceDescription = "Simple service,justforfun"

type Program struct{}

var (
	serviceIsRunning bool
	programIsRunning bool
	writingSync      sync.Mutex
)

func (p *Program) run() {
	logDirectoryCheck()
	go deleteOldLogFiles(48 * time.Hour)
	logInfo("RUN", "Program is running")
	router := httprouter.New()
	timer := sse.New()
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.GET("/", serveHomepage)

	router.POST("/get_time", getTime)

	router.Handler("GET", "/time", timer)
	go streamTime(timer)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Cannot start the server" + err.Error())
		os.Exit(-1)
	}
}

func (p *Program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	writingSync.Lock()
	serviceIsRunning = true
	writingSync.Unlock()
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	writingSync.Lock()
	serviceIsRunning = false
	writingSync.Unlock()
	for programIsRunning {
		fmt.Println(s.String() + " stopping")
		time.Sleep(time.Second * 1)
	}
	fmt.Println(s.String() + " stopped")
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &Program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service" + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot run the service" + err.Error())
	}
}
