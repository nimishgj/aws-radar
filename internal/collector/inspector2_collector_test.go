package collector

import "testing"

func TestCollector_inspector2_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewInspector2Collector(), "inspector2")
}

func TestCollector_inspector2_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewInspector2Collector(), true)
}
