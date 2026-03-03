package collector

import "testing"

func TestGlobalCollector_globalaccelerator_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewGlobalAcceleratorCollector(), "globalaccelerator")
}

func TestGlobalCollector_globalaccelerator_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewGlobalAcceleratorCollector(), true)
}
