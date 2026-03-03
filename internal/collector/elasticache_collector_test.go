package collector

import "testing"

func TestCollector_elasticache_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewElastiCacheCollector(), "elasticache")
}

func TestCollector_elasticache_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewElastiCacheCollector(), true)
}
