package collector

import "testing"

func TestCollector_secretsmanager_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewSecretsManagerCollector(), "secretsmanager")
}

func TestCollector_secretsmanager_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewSecretsManagerCollector(), true)
}
