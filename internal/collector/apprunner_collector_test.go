package collector

import "testing"

func TestCollector_apprunner_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAppRunnerCollector(), "apprunner")
}

func TestCollector_apprunner_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAppRunnerCollector(), true)
}
