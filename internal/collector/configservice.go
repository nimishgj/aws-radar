package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ConfigServiceCollector struct{}

func NewConfigServiceCollector() *ConfigServiceCollector { return &ConfigServiceCollector{} }

func (c *ConfigServiceCollector) Name() string { return "configservice" }

func (c *ConfigServiceCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := configservice.NewFromConfig(cfg)
	output, err := client.DescribeConfigurationRecorders(ctx, &configservice.DescribeConfigurationRecordersInput{})
	if err != nil {
		return err
	}

	counts := make(map[string]float64)
	for _, recorder := range output.ConfigurationRecorders {
		status := "UNKNOWN"
		if recorder.RecordingGroup != nil && recorder.RecordingGroup.AllSupported {
			status = "ALL_SUPPORTED"
		}
		counts[status]++
	}

	for status, count := range counts {
		metrics.ConfigRecorders.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
