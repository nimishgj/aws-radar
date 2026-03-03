package collector

import "testing"

func TestGlobalCollector_cost_explorer_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewCostCollector("daily"), "cost_explorer")
}

func TestGlobalCollector_cost_explorer_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewCostCollector("daily"), true)
}
