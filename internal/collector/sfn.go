package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SFNCollector struct{}

func NewSfnCollector() *SFNCollector {
	return &SFNCollector{}
}

func (c *SFNCollector) Name() string {
	return "sfn"
}

func (c *SFNCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := sfn.NewFromConfig(cfg)

	counts := make(map[string]float64)
	paginator := sfn.NewListStateMachinesPaginator(client, &sfn.ListStateMachinesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, sm := range page.StateMachines {
			smType := string(sm.Type)
			if smType == "" {
				smType = "unknown"
			}
			counts[smType]++
		}
	}

	for smType, count := range counts {
		metrics.SFNStateMachines.WithLabelValues(account, region, smType).Set(count)
	}

	log.Debug().
		Str("region", region).
		Int("state_machine_types", len(counts)).
		Msg("Step Functions collection completed")

	return nil
}
