package collector

import "testing"

func TestGlobalCollector_iam_Name(t *testing.T) {
	assertGlobalCollectorName(t, NewIAMCollector(), "iam")
}

func TestGlobalCollector_iam_ErrorContract(t *testing.T) {
	assertGlobalCollectorErrorContract(t, NewIAMCollector(), false)
}
