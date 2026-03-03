package collector

import "testing"

func TestCollector_athena_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAthenaCollector(), "athena")
}

func TestCollector_athena_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAthenaCollector(), true)
}
