package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type RDSCollector struct{}

func NewRDSCollector() *RDSCollector {
	return &RDSCollector{}
}

func (c *RDSCollector) Name() string {
	return "rds"
}

func (c *RDSCollector) Collect(ctx context.Context, cfg aws.Config, region string) error {
	client := rds.NewFromConfig(cfg)

	counts := make(map[string]float64)

	paginator := rds.NewDescribeDBInstancesPaginator(client, &rds.DescribeDBInstancesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, instance := range page.DBInstances {
			dbClass := aws.ToString(instance.DBInstanceClass)
			engine := aws.ToString(instance.Engine)
			multiAZ := strconv.FormatBool(instance.MultiAZ != nil && *instance.MultiAZ)
			status := aws.ToString(instance.DBInstanceStatus)
			if status == "" {
				status = "unknown"
			}

			key := dbClass + "|" + engine + "|" + multiAZ + "|" + status
			counts[key]++
		}
	}

	// Update metrics
	for key, count := range counts {
		parts := splitKey(key, 4)
		metrics.RDSInstances.WithLabelValues(
			region,
			parts[0], // db_instance_class
			parts[1], // engine
			parts[2], // multi_az
			parts[3], // status
		).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("instance_combinations", len(counts)).
		Msg("RDS collection completed")

	return nil
}
