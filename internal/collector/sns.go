package collector

import (
	"context"
	"strconv"
	"strings"

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

func (c *SNSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := sns.NewFromConfig(cfg)

	var totalSubscriptions float64
	typeCounts := make(map[string]float64) // fifo -> true/false

	paginator := sns.NewListTopicsPaginator(client, &sns.ListTopicsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, topic := range page.Topics {
			topicArn := aws.ToString(topic.TopicArn)

			// Determine FIFO vs standard from ARN
			isFifo := strings.HasSuffix(topicArn, ".fifo")
			typeCounts[strconv.FormatBool(isFifo)]++

			// Get topic attributes for subscription count
			attrs, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
				TopicArn: topic.TopicArn,
			})
			if err != nil {
				log.Warn().Err(err).Str("region", region).Str("topic", topicArn).Msg("Failed to get SNS topic attributes")
				continue
			}

			if subsStr, ok := attrs.Attributes["SubscriptionsConfirmed"]; ok {
				if subs, parseErr := strconv.ParseFloat(subsStr, 64); parseErr == nil {
					totalSubscriptions += subs
				}
			}
		}
	}

	// Total topics (preserve existing metric)
	var totalTopics float64
	for _, count := range typeCounts {
		totalTopics += count
	}
	metrics.SNSTopics.WithLabelValues(account, accountName, region).Set(totalTopics)

	// Topics by type
	for fifo, count := range typeCounts {
		metrics.SNSTopicsByType.WithLabelValues(account, accountName, region, fifo).Set(count)
	}

	// Total subscriptions
	metrics.SNSSubscriptions.WithLabelValues(account, accountName, region).Set(totalSubscriptions)

	log.Debug().
		Str("region", region).
		Float64("topic_count", totalTopics).
		Float64("subscriptions", totalSubscriptions).
		Int("topic_types", len(typeCounts)).
		Msg("SNS collection completed")

	return nil
}
