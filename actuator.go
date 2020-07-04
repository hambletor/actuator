package actuator

import (
	"log"
	"sync"
	"net/http"
)

//DefaultPort is the defalut port used to server actuator info if one is not given
// or is not a valid port number
const DefaultPort uint = 2121

// Health is the path for health check results
const Health string = "/actuator/health/"

// const history string = "/actuator/history/"

const info string = "/actuator/info/"

const environment string = "/actuator/envs/"

const metrics string = "/acutator/metreics/"

const up string = "UP"
const down string = "DOWN"

var (
	a    *Actuator
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
		a := &Actuator{Build: info}
		if check != nil {
			a.status.healthCheck = check
			log.Println("Health Check function set")
		}
		a.port = port
		if port <= 1024 || port > 65535 {
			a.port = DefaultPort
			log.Printf("port %d is not valid, using %d for Actuator", port, a.port)
		}
		// start http listener in seperate go routine
		go a.start()
		log.Println("Actuator initialized")
	})
	return a
}

func (a Actuator) start() {
	// init and start http server
	http.HandleFunc(Health,a.healthHandler())
	http.HandleFunc(info,a.infoHandler())
}

