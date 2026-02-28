package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type LambdaCollector struct{}

func NewLambdaCollector() *LambdaCollector {
	return &LambdaCollector{}
}

func (c *LambdaCollector) Name() string {
	return "lambda"
}

func (c *LambdaCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := lambda.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := lambda.NewListFunctionsPaginator(client, &lambda.ListFunctionsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, fn := range page.Functions {
			runtime := string(fn.Runtime)
			if runtime == "" {
				runtime = "unknown"
			}
			memorySize := strconv.FormatInt(int64(aws.ToInt32(fn.MemorySize)), 10)

			key := runtime + "|" + memorySize
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 2)
		metrics.LambdaFunctions.WithLabelValues(account, accountName, region,
			parts[0], // runtime
			parts[1], // memory_size
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("function_combinations", len(counts)).
		Msg("Lambda collection completed")

	return nil
}
