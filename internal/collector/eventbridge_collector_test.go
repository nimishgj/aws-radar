package collector

import "testing"

func TestCollector_eventbridge_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewEventBridgeCollector(), "eventbridge")
}

func TestCollector_eventbridge_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewEventBridgeCollector(), true)
}
