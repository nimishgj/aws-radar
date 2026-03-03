package collector

import "testing"

func TestCollector_securityhub_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSecurityHubCollector(), "securityhub")
}

func TestCollector_securityhub_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSecurityHubCollector(), true)
}
