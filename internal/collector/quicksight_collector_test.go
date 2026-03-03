package collector

import "testing"

func TestCollector_quicksight_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewQuickSightCollector(), "quicksight")
}

func TestCollector_quicksight_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewQuickSightCollector(), true)
}
