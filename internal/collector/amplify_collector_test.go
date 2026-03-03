package collector

import "testing"

func TestCollector_amplify_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAmplifyCollector(), "amplify")
}

func TestCollector_amplify_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAmplifyCollector(), true)
}
