package actuator

import (
	"testing"
)


func TestNewActuator(t *testing.T) {
	b := &BuildInfo{Name: "hambletor", Version: "version"}
	a := NewActuator(b, nil,8888)
	
	if a.Build.Name != "hambletor" {
		t.Errorf("expecting %s, got %s\n", "hambletor", a.Build.Name)
	}
}


func TestNewActuatorBadPort(t *testing.T) {
	b := &BuildInfo{Name: "hambletor", Version: "version"}
	a := NewActuator(b, nil,1)
	if a.port != DefaultPort{
		t.Errorf("expecting %d, got %d\n", DefaultPort, a.port)
	}
}

func TestActuatorHealthCehck(*testing.T) {

}
