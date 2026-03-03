package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type GlobalAcceleratorCollector struct{}

func NewGlobalAcceleratorCollector() *GlobalAcceleratorCollector {
	return &GlobalAcceleratorCollector{}
}

func (c *GlobalAcceleratorCollector) Name() string { return "globalaccelerator" }

func (c *GlobalAcceleratorCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	cfg.Region = "us-west-2"
	client := globalaccelerator.NewFromConfig(cfg)
	paginator := globalaccelerator.NewListAcceleratorsPaginator(client, &globalaccelerator.ListAcceleratorsInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, acc := range page.Accelerators {
			ipType := string(acc.IpAddressType)
			if ipType == "" {
				ipType = "UNKNOWN"
			}
			enabled := strconv.FormatBool(acc.Enabled != nil && *acc.Enabled)
			counts[ipType+"|"+enabled]++
		}
	}

	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.GlobalAccelerators.WithLabelValues(account, accountName, parts[0], parts[1]).Set(count)
	}
	return nil
}
