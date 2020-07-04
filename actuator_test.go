package actuator

import (
	"testing"
)

func TestNewActuator(t *testing.T) {
	b := &BuildInfo{Name: "hambletor", Version: "version"}
	a := NewActuator(b, nil)

	if a.Info.Name != "hambletor" {
		t.Errorf("expecting %s, got %s\n", "hambletor", a.Info.Name)
	}
}

func TestActuatorHealthCehck(*testing.T) {

}
