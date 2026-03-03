package collector

import "testing"

func TestCollector_backup_Name(t *testing.T) {
	assertRegionalCollectorName(t, NewBackupCollector(), "backup")
}

func TestCollector_backup_ErrorContract(t *testing.T) {
	assertRegionalCollectorErrorContract(t, NewBackupCollector(), true)
}
