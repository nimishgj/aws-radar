package collector

import "testing"

func TestCollector_neptune_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewNeptuneCollector(), "neptune")
}

func TestCollector_neptune_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewNeptuneCollector(), true)
}
