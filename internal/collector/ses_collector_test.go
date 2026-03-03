package collector

import "testing"

func TestCollector_ses_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSESCollector(), "ses")
}

func TestCollector_ses_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSESCollector(), true)
}
