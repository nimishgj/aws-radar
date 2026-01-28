package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type IAMCollector struct{}

func NewIAMCollector() *IAMCollector {
	return &IAMCollector{}
}

func (c *IAMCollector) Name() string {
	return "iam"
}

func (c *IAMCollector) Collect(ctx context.Context, cfg aws.Config) error {
	// IAM is a global service
	cfg.Region = "us-east-1"
	client := iam.NewFromConfig(cfg)

	// Collect users
	if err := c.collectUsers(ctx, client); err != nil {
		log.Warn().Err(err).Msg("Failed to collect IAM users")
	}

	// Collect roles
	if err := c.collectRoles(ctx, client); err != nil {
		log.Warn().Err(err).Msg("Failed to collect IAM roles")
	}

	return nil
}

func (c *IAMCollector) collectUsers(ctx context.Context, client *iam.Client) error {
	var count float64 = 0

	paginator := iam.NewListUsersPaginator(client, &iam.ListUsersInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		count += float64(len(page.Users))
	}

	metrics.IAMUsers.WithLabelValues().Set(count)

	log.Debug().
		Float64("user_count", count).
		Msg("IAM users collection completed")

	return nil
}

func (c *IAMCollector) collectRoles(ctx context.Context, client *iam.Client) error {
	var count float64 = 0

	paginator := iam.NewListRolesPaginator(client, &iam.ListRolesInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		count += float64(len(page.Roles))
	}

	metrics.IAMRoles.WithLabelValues().Set(count)

	log.Debug().
		Float64("role_count", count).
		Msg("IAM roles collection completed")

	return nil
}
