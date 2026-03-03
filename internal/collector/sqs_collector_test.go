package collector

import "testing"

func TestCollector_sqs_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSQSCollector(), "sqs")
}

func TestCollector_sqs_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSQSCollector(), true)
}
