package actuator

//DefaultPort is the defalut port used to server actuator info if one is not given
// or is not a valid port number
const DefaultPort uint = 2121

// Health is the path for health check results
const Health string = "health"

const history string = "history"

const info string = "info"

const environment string = "envs"

const metrics string = "metreics"

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
	//TODO look at thread saftey, this can be called by multiple sources simultaneously
	//TODO do we want to throttle this?
	a.status.previous = a.status.current
	if a.status.healthCheck != nil {
		a.status.current = "UP"
		if !a.status.healthCheck() {
			a.status.current = "DOWN"
		}
	}
	a.status.current = "Unknown, health check function not set"
	return a.status.current
}

//NewActuator creates a new actuator
func NewActuator(info *BuildInfo, check Check, port uint) *Actuator {
	a := &Actuator{Build: info}
	if check != nil {
		a.status.healthCheck = check
	}
	a.port = port
	if port < 1024 || port > 65535 {
		a.port = DefaultPort
	}
	return a
}
