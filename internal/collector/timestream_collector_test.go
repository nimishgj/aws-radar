package collector

import "testing"

func TestCollector_timestream_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewTimestreamCollector(), "timestream")
}

func TestCollector_timestream_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewTimestreamCollector(), true)
}
