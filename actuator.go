package actuator

import (
	"log"
	"net/http"
	"strconv"
	"sync"
)

//DefaultPort is the defalut port used to server actuator info if one is not given
// or is not a valid port number
const DefaultPort uint = 2121

const up string = "UP"
const down string = "DOWN"

var (
	singleton    *Actuator
	once sync.Once
	_hcm *sync.Mutex = &sync.Mutex{}
)

//Actuator is the base for all actuator information
type Actuator struct {
	Build  *BuildInfo `json:"build-info"`
	status Status
	port   uint
}

//BuildInfo contains the data around the build (expvar)
type BuildInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Date    string `json:"date"`
	Author  string `json:"author"`
}

//Check defines the function type
type Check func() bool

//Status shows the current status
type Status struct {
	//TODO think about ading some history of statuses and times when they changed.
	Current     string `json:"status"`
	previous    string
	healthCheck Check
}

//HealthCheck is the mechanism to execute the user defined Check
func (a *Actuator) HealthCheck() {
	_hcm.Lock()
	defer _hcm.Unlock()
	a.status.previous = a.status.Current
	if a.status.healthCheck != nil {
		a.status.Current = up
		if !a.status.healthCheck() {
			a.status.Current = down
		}
	}
	a.status.Current = "Unknown, health check function not set"
}

//NewActuator creates a new actuator or returns the singleton Actuator
func NewActuator(info *BuildInfo, check Check, port uint) *Actuator {
	once.Do(func() {
		singleton = &Actuator{Build: info, port:port}
		if check != nil {
			singleton.status.healthCheck = check
			log.Println("Health Check function set")
		}
		if port <= 1024 || port > 65535 {
			singleton.port = DefaultPort
			log.Printf("port %d is not valid, using %d for Actuator", port, singleton.port)
		}
		log.Println("Actuator initialized")
		singleton.start()
		log.Println("Actuator started")
	})
	return singleton
}

//start starts up the internal http listener
func (a Actuator) start() {
	// init and start http server
	http.HandleFunc(health, a.healthHandler())
	http.HandleFunc(info, a.infoHandler())
	http.HandleFunc(metrics, a.metricsHandler())
	http.HandleFunc(environment, a.envHandler())

	//spawn seperate listener
	go func() {
		log.Println("starting actuator listener")
		http.ListenAndServe(":"+strconv.Itoa(int(a.port)), nil)
	}()
}

