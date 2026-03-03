package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/nimishgj/aws-radar/internal/metrics"
)

type OrganizationsCollector struct{}

func NewOrganizationsCollector() *OrganizationsCollector { return &OrganizationsCollector{} }

func (c *OrganizationsCollector) Name() string { return "organizations" }

func (c *OrganizationsCollector) Collect(ctx context.Context, cfg aws.Config, account, accountName string) error {
	cfg.Region = "us-east-1"
	client := organizations.NewFromConfig(cfg)

	accountCounts := make(map[string]float64)
	accountPaginator := organizations.NewListAccountsPaginator(client, &organizations.ListAccountsInput{})
	for accountPaginator.HasMorePages() {
		page, err := accountPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, acc := range page.Accounts {
			state := string(acc.State)
			if state == "" {
				state = "UNKNOWN"
			}
			accountCounts[state]++
		}
	}
	for state, count := range accountCounts {
		metrics.OrganizationsAccounts.WithLabelValues(account, accountName, state).Set(count)
	}

	var ouCount float64
	rootPaginator := organizations.NewListRootsPaginator(client, &organizations.ListRootsInput{})
	for rootPaginator.HasMorePages() {
		rootPage, err := rootPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, root := range rootPage.Roots {
			ouPaginator := organizations.NewListOrganizationalUnitsForParentPaginator(client, &organizations.ListOrganizationalUnitsForParentInput{ParentId: root.Id})
			for ouPaginator.HasMorePages() {
				ouPage, err := ouPaginator.NextPage(ctx)
				if err != nil {
					return err
				}
				ouCount += float64(len(ouPage.OrganizationalUnits))
			}
		}
	}
	metrics.OrganizationsOrganizationalUnits.WithLabelValues(account, accountName).Set(ouCount)

	return nil
}
