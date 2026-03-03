package collector

import "testing"

func TestCollector_guardduty_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewGuardDutyCollector(), "guardduty")
}

func TestCollector_guardduty_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewGuardDutyCollector(), true)
}
