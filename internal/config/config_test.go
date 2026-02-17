package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadDefaults(t *testing.T) {
	viper.Reset()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Fatalf("expected default port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Collection.Interval.String() != "1m0s" {
		t.Fatalf("expected default interval 60s, got %s", cfg.Collection.Interval)
	}
	if len(cfg.AWS.Regions) != 1 || cfg.AWS.Regions[0] != "us-east-1" {
		t.Fatalf("expected default region us-east-1, got %v", cfg.AWS.Regions)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("expected default logging level info, got %s", cfg.Logging.Level)
	}
}
