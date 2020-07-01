package actuator

import (
	"log"
	"sync"
)

//DefaultPort is the defalut port used to server actuator info if one is not given
// or is not a valid port number
const DefaultPort uint = 2121

// Health is the path for health check results
const Health string = "health"

const history string = "history"

const info string = "info"

const environment string = "envs"

const metrics string = "metreics"

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
	status status
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

type status struct {
	//TODO think about ading some history of statuses and times when they changed.
	current     string
	previous    string
	healthCheck Check
}

//HealthCheck is the mechanism to execute the user defined Check
func (a *Actuator) HealthCheck() string {
	_hcm.Lock()
	defer _hcm.Unlock()
	a.status.previous = a.status.current
	if a.status.healthCheck != nil {
		a.status.current = up
		if !a.status.healthCheck() {
			a.status.current = down
		}
	}
	a.status.current = "Unknown, health check function not set"
	return a.status.current
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
		log.Println("Actuator initialized")
	})
	return a
}
