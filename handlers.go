package actuator

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	health      string = "/actuator/health/"
	info        string = "/actuator/info/"
	environment string = "/actuator/envs/"
	metrics     string = "/actuator/metrics/"
)

//HTTP Handlers for Actuator endpoints

func (a Actuator) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		a.HealthCheck()
		w.Write(marshal(a.status))
	}
}

func (a Actuator) envHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(envs())
	}
}

func (a Actuator) metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(a.stats())
	}
}

func (a Actuator) infoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(marshal(a.Info))
	}
}

//start starts up the internal http listener
func (a Actuator) addHandlers() {
	http.HandleFunc(health, a.healthHandler())
	http.HandleFunc(info, a.infoHandler())
	http.HandleFunc(metrics, a.metricsHandler())
	http.HandleFunc(environment, a.envHandler())
}

func (a Actuator) stats() []byte {
	runtime.ReadMemStats(a.Metrics.Mem)
	a.Metrics.Routines = runtime.NumGoroutine()
	return marshal(a.Metrics)
}

//Env holds envioronment variables
type Env struct {
	Variables []string `json:"envs"`
}

// always get the most current environment variables
func envs() []byte {
	var e Env
	e.Variables = os.Environ()
	return marshal(e)
}

func marshal(o interface{}) []byte {
	b, err := json.Marshal(o)
	if err != nil {
		log.Printf("issue marshaling Acutator Build Info: %s", err.Error())
		b = []byte("unable to mashal data")
	}
	return b
}
