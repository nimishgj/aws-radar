package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/controltower"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type ControlTowerCollector struct{}

func NewControlTowerCollector() *ControlTowerCollector { return &ControlTowerCollector{} }

func (c *ControlTowerCollector) Name() string { return "controltower" }

func (c *ControlTowerCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := controltower.NewFromConfig(cfg)
	paginator := controltower.NewListLandingZonesPaginator(client, &controltower.ListLandingZonesInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, lz := range page.LandingZones {
			status := "UNKNOWN"
			if lz.Arn != nil {
				detail, err := client.GetLandingZone(ctx, &controltower.GetLandingZoneInput{LandingZoneIdentifier: lz.Arn})
				if err == nil && detail.LandingZone != nil {
					status = string(detail.LandingZone.Status)
				}
			}
			if status == "" {
				status = "UNKNOWN"
			}
			counts[status]++
		}
	}

	for status, count := range counts {
		metrics.ControlTowerLandingZones.WithLabelValues(account, accountName, region, status).Set(count)
	}
	return nil
}
