package collector

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SSMCollector struct{}

func NewSSMCollector() *SSMCollector {
	return &SSMCollector{}
}

func (c *SSMCollector) Name() string {
	return "ssm"
}

func (c *SSMCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := ssm.NewFromConfig(cfg)

	// Parameters (existing)
	paramCounts := make(map[string]float64)
	paramPaginator := ssm.NewDescribeParametersPaginator(client, &ssm.DescribeParametersInput{})
	for paramPaginator.HasMorePages() {
		page, err := paramPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, param := range page.Parameters {
			paramType := string(param.Type)
			if paramType == "" {
				paramType = "unknown"
			}
			paramCounts[paramType]++
		}
	}
	for paramType, count := range paramCounts {
		metrics.SSMParameters.WithLabelValues(account, accountName, region, paramType).Set(count)
	}

	// Documents
	docCounts := make(map[string]float64)
	for _, owner := range []string{"Self", "Amazon", "ThirdParty"} {
		docPaginator := ssm.NewListDocumentsPaginator(client, &ssm.ListDocumentsInput{
			Filters: []ssmTypes.DocumentKeyValuesFilter{
				{Key: aws.String("Owner"), Values: []string{owner}},
			},
		})
		for docPaginator.HasMorePages() {
			page, err := docPaginator.NextPage(ctx)
			if err != nil {
				log.Warn().Err(err).Str("region", region).Str("owner", owner).Msg("Failed to list SSM documents")
				break
			}
			docCounts[owner] += float64(len(page.DocumentIdentifiers))
		}
	}
	for owner, count := range docCounts {
		metrics.SSMDocuments.WithLabelValues(account, accountName, region, owner).Set(count)
	}

	// Maintenance Windows
	mwCounts := make(map[string]float64)
	var mwNextToken *string
	for {
		output, err := client.DescribeMaintenanceWindows(ctx, &ssm.DescribeMaintenanceWindowsInput{
			NextToken: mwNextToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to describe SSM maintenance windows")
			break
		}
		for _, mw := range output.WindowIdentities {
			enabled := strconv.FormatBool(mw.Enabled)
			mwCounts[enabled]++
		}
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		mwNextToken = output.NextToken
	}
	for enabled, count := range mwCounts {
		metrics.SSMMaintenanceWindows.WithLabelValues(account, accountName, region, enabled).Set(count)
	}

	// Associations
	var assocCount float64
	assocPaginator := ssm.NewListAssociationsPaginator(client, &ssm.ListAssociationsInput{})
	for assocPaginator.HasMorePages() {
		page, err := assocPaginator.NextPage(ctx)
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to list SSM associations")
			break
		}
		assocCount += float64(len(page.Associations))
	}
	metrics.SSMAssociations.WithLabelValues(account, accountName, region).Set(assocCount)

	// Patch Baselines
	var patchCount float64
	var patchNextToken *string
	for {
		output, err := client.DescribePatchBaselines(ctx, &ssm.DescribePatchBaselinesInput{
			NextToken: patchNextToken,
		})
		if err != nil {
			log.Warn().Err(err).Str("region", region).Msg("Failed to describe SSM patch baselines")
			break
		}
		patchCount += float64(len(output.BaselineIdentities))
		if output.NextToken == nil || *output.NextToken == "" {
			break
		}
		patchNextToken = output.NextToken
	}
	metrics.SSMPatchBaselines.WithLabelValues(account, accountName, region).Set(patchCount)

	log.Debug().
		Str("region", region).
		Int("parameter_types", len(paramCounts)).
		Int("document_owners", len(docCounts)).
		Int("maintenance_window_states", len(mwCounts)).
		Float64("associations", assocCount).
		Float64("patch_baselines", patchCount).
		Msg("SSM collection completed")

	return nil
}
