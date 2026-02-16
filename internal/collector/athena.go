package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type AthenaCollector struct{}

func NewAthenaCollector() *AthenaCollector {
	return &AthenaCollector{}
}

func (c *AthenaCollector) Name() string {
	return "athena"
}

func (c *AthenaCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := athena.NewFromConfig(cfg)

	var count float64
	paginator := athena.NewListWorkGroupsPaginator(client, &athena.ListWorkGroupsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.WorkGroups))
	}

	metrics.AthenaWorkgroups.WithLabelValues(region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("workgroup_count", count).
		Msg("Athena workgroup collection completed")

	return nil
}
