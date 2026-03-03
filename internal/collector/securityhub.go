package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type SecurityHubCollector struct{}

func NewSecurityHubCollector() *SecurityHubCollector {
	return &SecurityHubCollector{}
}

func (c *SecurityHubCollector) Name() string {
	return "securityhub"
}

func (c *SecurityHubCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := securityhub.NewFromConfig(cfg)
	paginator := securityhub.NewGetEnabledStandardsPaginator(client, &securityhub.GetEnabledStandardsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.StandardsSubscriptions))
	}

	metrics.SecurityHubStandards.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
