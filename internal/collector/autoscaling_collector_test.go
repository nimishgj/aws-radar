package collector

import "testing"

func TestCollector_autoscaling_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewAutoScalingCollector(), "autoscaling")
}

func TestCollector_autoscaling_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewAutoScalingCollector(), true)
}
