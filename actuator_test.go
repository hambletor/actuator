package actuator

import "testing"

// type test info
// var tests = []testpair{
// 	{ []float64{1,2}, 1.5 },
// 	{ []float64{1,1,1,1,1,1}, 1 },
// 	{ []float64{-1,1}, 0 },
//   }

func TestNewActuator(t *testing.T) {
	b := &actBuildInfo{Name: "hambletor", Version: "version"}
	a := NewActuator(b, nil)
	if a.Build.Name != "hambletor" {
		t.Errorf("expecting %s, got %s\n", "hambletor", a.Build.Name)
	}
}

func TestActuatorHealthCehck(*testing.T) {

}
