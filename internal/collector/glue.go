package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type GlueCollector struct{}

func NewGlueCollector() *GlueCollector {
	return &GlueCollector{}
}

func (c *GlueCollector) Name() string {
	return "glue"
}

func (c *GlueCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := glue.NewFromConfig(cfg)

	var count float64
	paginator := glue.NewGetJobsPaginator(client, &glue.GetJobsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		count += float64(len(page.Jobs))
	}

	metrics.GlueJobs.WithLabelValues(account, accountName, region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("job_count", count).
		Msg("Glue job collection completed")

	return nil
}
