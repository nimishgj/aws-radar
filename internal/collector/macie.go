package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/macie2"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type MacieCollector struct{}

func NewMacieCollector() *MacieCollector {
	return &MacieCollector{}
}

func (c *MacieCollector) Name() string {
	return "macie"
}

func (c *MacieCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := macie2.NewFromConfig(cfg)
	paginator := macie2.NewListClassificationJobsPaginator(client, &macie2.ListClassificationJobsInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Items))
	}

	metrics.MacieClassificationJobs.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
