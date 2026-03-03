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
	if cfg.CostExplorer.Enabled {
		t.Fatalf("expected cost_explorer.enabled default false, got true")
	}
	if cfg.CostExplorer.Frequency != "daily" {
		t.Fatalf("expected cost_explorer.frequency default daily, got %s", cfg.CostExplorer.Frequency)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("expected default logging level info, got %s", cfg.Logging.Level)
	}

	requiredCollectors := []string{
		"mq", "ses", "cloudformation", "documentdb", "neptune", "memorydb",
		"timestream", "fsx", "backup", "kinesis", "firehose", "kinesisanalytics",
		"emr", "elasticbeanstalk", "kms", "cloudtrail", "batch",
		"kinesisvideo", "opensearchserverless", "s3control",
		"codebuild", "codepipeline", "codedeploy",
		"apprunner", "transfer", "msk", "redshift",
		"guardduty", "securityhub", "inspector2", "macie", "waf",
		"route53resolver", "configservice", "organizations", "controltower",
		"ecrpublic", "directconnect", "bedrock", "sagemaker", "quicksight",
		"workspaces", "appstream", "connect", "amplify", "globalaccelerator",
		"datasync", "dms", "shield",
	}
	for _, collector := range requiredCollectors {
		found := false
		for _, enabled := range cfg.Collectors {
			if enabled == collector {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected default collectors to include %q", collector)
		}
	}
}
