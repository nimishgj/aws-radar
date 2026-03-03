package collector

import "testing"

func TestCollector_transfer_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewTransferCollector(), "transfer")
}

func TestCollector_transfer_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewTransferCollector(), true)
}
