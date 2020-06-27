package github.com/hambletor/actuator

//Actuator is the base for all actuator information
type Actuator struct {
	Build  BuildInfo `json:"build-info"`
	status status
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
	status      string
	healthCheck Check
}

func (a *Actuator) checkHealth() {
	if a.status.healthCheck != nil {
		a.status.status = "UP"
		if !a.status.healthCheck() {
			a.status.status = "DOWN"
		}
	}
	a.status.status = "Unknown, health check function not set"
}

func (a Actuator) health() string {
	return a.status.status
}

//NewActuator creates a new actuator
func NewActuator(info BuildInfo, check Check) *Actuator {
	a := &Actuator{Build: info}
	if check != nil {
		a.status.healthCheck = check
	}
	return a
}
