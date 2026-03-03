package collector

import "testing"

func TestCollector_kinesisanalytics_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewKinesisAnalyticsCollector(), "kinesisanalytics")
}

func TestCollector_kinesisanalytics_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewKinesisAnalyticsCollector(), true)
}
