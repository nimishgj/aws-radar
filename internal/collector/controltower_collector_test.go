package collector

import "testing"

func TestCollector_controltower_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewControlTowerCollector(), "controltower")
}

func TestCollector_controltower_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewControlTowerCollector(), true)
}
