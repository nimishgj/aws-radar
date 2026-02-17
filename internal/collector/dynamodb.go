package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type DynamoDBCollector struct{}

func NewDynamoDBCollector() *DynamoDBCollector {
	return &DynamoDBCollector{}
}

func (c *DynamoDBCollector) Name() string {
	return "dynamodb"
}

func (c *DynamoDBCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := dynamodb.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := dynamodb.NewListTablesPaginator(client, &dynamodb.ListTablesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, tableName := range page.TableNames {
			// Get table details
			descOutput, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: aws.String(tableName),
			})
			if err != nil {
				log.Warn().
					Err(err).
					Str("table", tableName).
					Msg("Failed to describe DynamoDB table")
				continue
			}

			billingMode := "PROVISIONED"
			if descOutput.Table.BillingModeSummary != nil {
				billingMode = string(descOutput.Table.BillingModeSummary.BillingMode)
			}
			counts[billingMode]++
		}
	}

	// Update metrics
	for billingMode, count := range counts {
		metrics.DynamoDBTables.WithLabelValues(account, region, billingMode).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("billing_modes", len(counts)).
		Msg("DynamoDB collection completed")

	return nil
}
