package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type AppRunnerCollector struct{}

func NewAppRunnerCollector() *AppRunnerCollector {
	return &AppRunnerCollector{}
}

func (c *AppRunnerCollector) Name() string {
	return "apprunner"
}

func (c *AppRunnerCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := apprunner.NewFromConfig(cfg)
	paginator := apprunner.NewListServicesPaginator(client, &apprunner.ListServicesInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.ServiceSummaryList))
	}

	metrics.AppRunnerServices.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
