package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type Inspector2Collector struct{}

func NewInspector2Collector() *Inspector2Collector {
	return &Inspector2Collector{}
}

func (c *Inspector2Collector) Name() string {
	return "inspector2"
}

func (c *Inspector2Collector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := inspector2.NewFromConfig(cfg)
	paginator := inspector2.NewListCoveragePaginator(client, &inspector2.ListCoverageInput{})

	var count float64
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.CoveredResources))
	}

	metrics.InspectorCoveredResources.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
