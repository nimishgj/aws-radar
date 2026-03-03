package collector

import "testing"

func TestCollector_ssm_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSSMCollector(), "ssm")
}

func TestCollector_ssm_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSSMCollector(), true)
}
