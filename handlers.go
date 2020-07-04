package actuator

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
)

//HTTP Handlers for each endpoint

func (a Actuator) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		a.HealthCheck()
		w.Write(marshal(a.status))
	}
}

func (a Actuator)envHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(envs())
	}
}

func (a Actuator) metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(stats())
	}
}

func (a Actuator) infoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(marshal(a.Build))
	}
}

func marshal(o interface{}) []byte {
	b, err := json.Marshal(a.Build)
	if err != nil {
		log.Printf("issue marshaling Acutator Build Info: %s", err.Error())
		b = []byte("unable to mashal data")
	}
	return b
}

//Monitor holds the metrics for the app
// Thanks https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/
type Monitor struct {
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"totalAlloc"`
	Sys          uint64 `json:"sys"`
	Mallocs      uint64 `json"mallocs"`
	Frees        uint64 `json:"frees"`
	LiveObjects  uint64 `json:"liveObjects"`
	PauseTotalNs uint64 `json:"pauseTotalNs"`
	NumGC        uint32 `json:"numGC"`
	NumGoroutine int    `json:"numGoroutines"`
}

func stats() []byte {
	var m Monitor
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	// Number of goroutines
	m.NumGoroutine = runtime.NumGoroutine()

	// Misc memory stats
	m.Alloc = rtm.Alloc
	m.TotalAlloc = rtm.TotalAlloc
	m.Sys = rtm.Sys
	m.Mallocs = rtm.Mallocs
	m.Frees = rtm.Frees

	// Live objects = Mallocs - Frees
	m.LiveObjects = m.Mallocs - m.Frees

	// GC Stats
	m.PauseTotalNs = rtm.PauseTotalNs
	m.NumGC = rtm.NumGC

	// Just encode to json and print
	return marshal(m)
}

//Env holds envioronment variables
type Env struct {
	Variables []string `json:"envs"`
}

func envs() []byte {
	var e Env
	e.Variables = os.Environ()
	return marshal(e)
}
