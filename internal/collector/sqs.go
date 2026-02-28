package collector

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	paginator := sqs.NewListQueuesPaginator(client, &sqs.ListQueuesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, queueUrl := range page.QueueUrls {
			// Determine queue type (FIFO vs Standard)
			queueType := "standard"
			if strings.HasSuffix(queueUrl, ".fifo") {
				queueType = "fifo"
			}
			counts[queueType]++
		}
	}

	// Update metrics
	for queueType, count := range counts {
		metrics.SQSQueues.WithLabelValues(account, accountName, region, queueType).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("queue_types", len(counts)).
		Msg("SQS collection completed")

	return nil
}
