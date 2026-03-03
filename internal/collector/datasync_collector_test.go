package collector

import "testing"

func TestCollector_datasync_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewDataSyncCollector(), "datasync")
}

func TestCollector_datasync_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewDataSyncCollector(), true)
}
