package collector

import "testing"

func TestCollector_batch_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewBatchCollector(), "batch")
}

func TestCollector_batch_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewBatchCollector(), true)
}
