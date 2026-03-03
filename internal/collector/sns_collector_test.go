package collector

import "testing"

func TestCollector_sns_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSNSCollector(), "sns")
}

func TestCollector_sns_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSNSCollector(), true)
}
