package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type EventBridgeCollector struct{}

func NewEventBridgeCollector() *EventBridgeCollector {
	return &EventBridgeCollector{}
}

func (c *EventBridgeCollector) Name() string {
	return "eventbridge"
}

func (c *EventBridgeCollector) Collect(ctx context.Context, cfg aws.Config, region, account string) error {
	client := eventbridge.NewFromConfig(cfg)

	totalBuses := 0
	var busToken *string
	for {
		page, err := client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{
			NextToken: busToken,
		})
		if err != nil {
			return err
		}

		for _, bus := range page.EventBuses {
			busName := aws.ToString(bus.Name)
			if busName == "" {
				busName = "default"
			}

			var ruleCount float64
			var ruleToken *string
			for {
				rulePage, err := client.ListRules(ctx, &eventbridge.ListRulesInput{
					EventBusName: aws.String(busName),
					NextToken:    ruleToken,
				})
				if err != nil {
					return err
				}
				ruleCount += float64(len(rulePage.Rules))
				if rulePage.NextToken == nil || len(*rulePage.NextToken) == 0 {
					break
				}
				ruleToken = rulePage.NextToken
			}

			metrics.EventBridgeRules.WithLabelValues(account, region, busName).Set(ruleCount)
			totalBuses++
		}
		if page.NextToken == nil || len(*page.NextToken) == 0 {
			break
		}
		busToken = page.NextToken
	}

	log.Debug().
		Str("region", region).
		Int("event_buses", totalBuses).
		Msg("EventBridge collection completed")

	return nil
}
