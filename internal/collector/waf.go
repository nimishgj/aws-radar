package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafTypes "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type WAFCollector struct{}

func NewWAFCollector() *WAFCollector {
	return &WAFCollector{}
}

func (c *WAFCollector) Name() string {
	return "waf"
}

func (c *WAFCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := wafv2.NewFromConfig(cfg)
	var count float64
	var marker *string
	for {
		page, err := client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
			Scope:      wafTypes.ScopeRegional,
			NextMarker: marker,
		})
		if err != nil {
			return err
		}
		count += float64(len(page.WebACLs))
		if page.NextMarker == nil || *page.NextMarker == "" {
			break
		}
		marker = page.NextMarker
	}

	metrics.WAFWebACLs.WithLabelValues(account, accountName, region).Set(count)
	return nil
}
