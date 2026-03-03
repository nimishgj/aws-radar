package collector

import "testing"

func TestCollector_mq_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewMQCollector(), "mq")
}

func TestCollector_mq_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewMQCollector(), true)
}
