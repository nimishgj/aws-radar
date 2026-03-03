package collector

import "testing"

func TestCollector_firehose_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewFirehoseCollector(), "firehose")
}

func TestCollector_firehose_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewFirehoseCollector(), true)
}
