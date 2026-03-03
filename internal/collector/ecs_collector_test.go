package collector

import "testing"

func TestCollector_ecs_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewECSCollector(), "ecs")
}

func TestCollector_ecs_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewECSCollector(), true)
}
