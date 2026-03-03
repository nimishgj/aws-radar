package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type GuardDutyCollector struct{}

func NewGuardDutyCollector() *GuardDutyCollector {
	return &GuardDutyCollector{}
}

func (c *GuardDutyCollector) Name() string {
	return "guardduty"
}

func (c *GuardDutyCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := guardduty.NewFromConfig(cfg)
	paginator := guardduty.NewListDetectorsPaginator(client, &guardduty.ListDetectorsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.DetectorIds))
	}

	metrics.GuardDutyDetectors.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
