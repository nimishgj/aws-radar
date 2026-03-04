package collector

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SQSCollector struct{}

func NewSQSCollector() *SQSCollector {
	return &SQSCollector{}
}

func (c *SQSCollector) Name() string {
	return "sqs"
}

func (c *SQSCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := sqs.NewFromConfig(cfg)

	counts := make(map[string]float64)
	messageCounts := make(map[string]float64)
	var dlqCount float64

	paginator := sqs.NewListQueuesPaginator(client, &sqs.ListQueuesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, queueUrl := range page.QueueUrls {
			queueType := "standard"
			if strings.HasSuffix(queueUrl, ".fifo") {
				queueType = "fifo"
			}
			counts[queueType]++

			// Get queue attributes for message count and DLQ config
			attrs, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
				QueueUrl: aws.String(queueUrl),
				AttributeNames: []sqsTypes.QueueAttributeName{
					sqsTypes.QueueAttributeNameApproximateNumberOfMessages,
					sqsTypes.QueueAttributeNameRedrivePolicy,
				},
			})
			if err != nil {
				log.Warn().Err(err).Str("region", region).Str("queue", queueUrl).Msg("Failed to get SQS queue attributes")
				continue
			}

			if msgStr, ok := attrs.Attributes["ApproximateNumberOfMessages"]; ok {
				if msgs, parseErr := strconv.ParseFloat(msgStr, 64); parseErr == nil {
					messageCounts[queueType] += msgs
				}
			}

			if _, ok := attrs.Attributes["RedrivePolicy"]; ok {
				dlqCount++
			}
		}
	}

	for queueType, count := range counts {
		metrics.SQSQueues.WithLabelValues(account, accountName, region, queueType).Set(count)
	}

	for queueType, msgs := range messageCounts {
		metrics.SQSMessages.WithLabelValues(account, accountName, region, queueType).Set(msgs)
	}

	metrics.SQSQueuesWithDLQ.WithLabelValues(account, accountName, region).Set(dlqCount)

	log.Debug().
		Str("region", region).
		Int("queue_types", len(counts)).
		Float64("dlq_queues", dlqCount).
		Msg("SQS collection completed")

	return nil
}
