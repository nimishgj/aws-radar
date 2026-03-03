package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type QuickSightCollector struct{}

func NewQuickSightCollector() *QuickSightCollector { return &QuickSightCollector{} }

func (c *QuickSightCollector) Name() string { return "quicksight" }

func (c *QuickSightCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := quicksight.NewFromConfig(cfg)
	paginator := quicksight.NewListDashboardsPaginator(client, &quicksight.ListDashboardsInput{AwsAccountId: aws.String(account)})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.DashboardSummaryList))
	}
	metrics.QuickSightDashboards.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
