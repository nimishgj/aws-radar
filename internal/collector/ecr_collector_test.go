package collector

import "testing"

func TestCollector_ecr_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewECRCollector(), "ecr")
}

func TestCollector_ecr_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewECRCollector(), true)
}
