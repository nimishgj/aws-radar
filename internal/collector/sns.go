package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SNSCollector struct{}

func NewSNSCollector() *SNSCollector {
	return &SNSCollector{}
}

func (c *SNSCollector) Name() string {
	return "sns"
}

func (c *SNSCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := sns.NewFromConfig(cfg)

	var count float64 = 0

	paginator := sns.NewListTopicsPaginator(client, &sns.ListTopicsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		count += float64(len(page.Topics))
	}

	metrics.SNSTopics.WithLabelValues(region).Set(count)

	log.Debug().
		Str("region", region).
		Float64("topic_count", count).
		Msg("SNS collection completed")

	return nil
}
