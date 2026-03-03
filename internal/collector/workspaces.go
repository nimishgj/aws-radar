package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type WorkSpacesCollector struct{}

func NewWorkSpacesCollector() *WorkSpacesCollector { return &WorkSpacesCollector{} }

func (c *WorkSpacesCollector) Name() string { return "workspaces" }

func (c *WorkSpacesCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := workspaces.NewFromConfig(cfg)
	paginator := workspaces.NewDescribeWorkspacesPaginator(client, &workspaces.DescribeWorkspacesInput{})

	counts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, ws := range page.Workspaces {
			state := string(ws.State)
			if state == "" {
				state = "UNKNOWN"
			}
			counts[state]++
		}
	}

	for state, count := range counts {
		metrics.WorkSpacesInstances.WithLabelValues(account, accountName, region, state).Set(count)
	}
	return nil
}
