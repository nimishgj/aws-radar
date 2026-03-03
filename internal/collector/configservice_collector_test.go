package collector

import "testing"

func TestCollector_configservice_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewConfigServiceCollector(), "configservice")
}

func TestCollector_configservice_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewConfigServiceCollector(), true)
}
