package actuator

import (
	"log"
	"runtime"
	"sync"
)


const up string = "UP"
const down string = "DOWN"

var (
	singleton *Actuator
	once      sync.Once
	_hcm      *sync.Mutex = &sync.Mutex{}
)

//Actuator is the base for all actuator information
type Actuator struct {
	Info     *BuildInfo `json:"build-info"`
	status   Status
	Metrics  Metrics
}

//BuildInfo contains the data around the build (expvar)
type BuildInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Date    string `json:"date"`
	Author  string `json:"author"`
}

//Metrics captures metrics from the runtime
type Metrics struct {
	Mem *runtime.MemStats `json:"memory"`
	Routines int `json:"goRoutines"`
}

//Check defines the function type
type Check func() bool

//Status shows the current status
type Status struct {
	Current string `json:"status"`
	check   Check
}

//HealthCheck is the mechanism to execute the user defined Check
func (a *Actuator) HealthCheck() {
	_hcm.Lock()
	defer _hcm.Unlock()
	if a.status.check != nil {
		a.status.Current = up
		if !a.status.check() {
			a.status.Current = down
		}
	}
	a.status.Current = "Unknown, health check function not set"
}

//NewActuator creates a new actuator or returns the singleton Actuator
func NewActuator(info *BuildInfo, check Check) *Actuator {
	once.Do(func() {
		singleton = &Actuator{Info: info, Metrics: Metrics{Mem:&runtime.MemStats{}}}
		if check != nil {
			singleton.status.check = check
			log.Println("Health Check function set")
		}
		log.Println("Actuator adding handlers")
		singleton.addHandlers()
		log.Println("Actuator initialized")
	})
	return singleton
}
