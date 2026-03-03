package collector

import "testing"

func TestCollector_sfn_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSfnCollector(), "sfn")
}

func TestCollector_sfn_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSfnCollector(), true)
}
