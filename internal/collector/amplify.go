package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type AmplifyCollector struct{}

func NewAmplifyCollector() *AmplifyCollector { return &AmplifyCollector{} }

func (c *AmplifyCollector) Name() string { return "amplify" }

func (c *AmplifyCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := amplify.NewFromConfig(cfg)
	paginator := amplify.NewListAppsPaginator(client, &amplify.ListAppsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Apps))
	}
	metrics.AmplifyApps.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
