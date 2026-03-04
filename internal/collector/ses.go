package collector

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2Types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

type SESCollector struct{}

func NewSESCollector() *SESCollector {
	return &SESCollector{}
}

func (c *SESCollector) Name() string {
	return "ses"
}

func (c *SESCollector) Collect(ctx context.Context, cfg aws.Config, region, account, accountName string) error {
	client := sesv2.NewFromConfig(cfg)
	paginator := sesv2.NewListEmailIdentitiesPaginator(client, &sesv2.ListEmailIdentitiesInput{})

	var count float64
	verificationCounts := make(map[string]float64)
	authStatusCounts := make(map[string]float64)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, identity := range page.EmailIdentities {
			count++
			verificationStatus := string(identity.VerificationStatus)
			if verificationStatus == "" {
				verificationStatus = "UNKNOWN"
			}
			verificationCounts[verificationStatus]++

			identityName := aws.ToString(identity.IdentityName)
			if identityName == "" {
				continue
			}
			detail, err := client.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
				EmailIdentity: aws.String(identityName),
			})
			if err != nil {
				log.Warn().Err(err).Str("region", region).Str("identity", identityName).Msg("Failed to get SES email identity")
				continue
			}
			dkimStatus := "UNKNOWN"
			if detail.DkimAttributes != nil {
				dkimStatus = string(detail.DkimAttributes.Status)
				if dkimStatus == "" {
					dkimStatus = "UNKNOWN"
				}
			}
			mailFromStatus := "NOT_CONFIGURED"
			if detail.MailFromAttributes != nil {
				mailFromStatus = string(detail.MailFromAttributes.MailFromDomainStatus)
				if mailFromStatus == "" {
					mailFromStatus = "UNKNOWN"
				}
			}
			spfStatus := "NOT_CONFIGURED"
			if detail.MailFromAttributes != nil {
				spfStatus = mailFromStatus
			}
			authKey := dkimStatus + "|" + spfStatus + "|" + mailFromStatus
			authStatusCounts[authKey]++
		}
	}

	metrics.SESIdentities.WithLabelValues(account, accountName, region).Set(count)
	for verificationStatus, c := range verificationCounts {
		metrics.SESIdentitiesByVerificationStatus.WithLabelValues(account, accountName, region, verificationStatus).Set(c)
	}
	for key, c := range authStatusCounts {
		parts := splitKey(key, 3)
		metrics.SESIdentityAuthStatus.WithLabelValues(account, accountName, region, parts[0], parts[1], parts[2]).Set(c)
	}

	var configSetCount float64
	eventDestinationCounts := make(map[string]float64)
	configSetByPoolCounts := make(map[string]float64)
	configSetPaginator := sesv2.NewListConfigurationSetsPaginator(client, &sesv2.ListConfigurationSetsInput{})
	for configSetPaginator.HasMorePages() {
		page, err := configSetPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, configSetName := range page.ConfigurationSets {
			configSetCount++
			eventDestinations, eventDestErr := client.GetConfigurationSetEventDestinations(ctx, &sesv2.GetConfigurationSetEventDestinationsInput{
				ConfigurationSetName: aws.String(configSetName),
			})
			if eventDestErr != nil {
				log.Warn().Err(eventDestErr).Str("region", region).Str("config_set", configSetName).Msg("Failed to get SES config set event destinations")
				continue
			}
			for _, destination := range eventDestinations.EventDestinations {
				for _, destinationType := range eventDestinationTypes(destination) {
					eventDestinationCounts[destinationType]++
				}
			}

			configSetDetail, configSetErr := client.GetConfigurationSet(ctx, &sesv2.GetConfigurationSetInput{
				ConfigurationSetName: aws.String(configSetName),
			})
			if configSetErr != nil {
				log.Warn().Err(configSetErr).Str("region", region).Str("config_set", configSetName).Msg("Failed to get SES configuration set")
				continue
			}
			sendingPool := "none"
			if configSetDetail.DeliveryOptions != nil {
				pool := aws.ToString(configSetDetail.DeliveryOptions.SendingPoolName)
				if pool != "" {
					sendingPool = pool
				}
			}
			configSetByPoolCounts[sendingPool]++
		}
	}
	metrics.SESConfigSets.WithLabelValues(account, accountName, region).Set(configSetCount)
	for destinationType, c := range eventDestinationCounts {
		metrics.SESConfigSetEventDestinations.WithLabelValues(account, accountName, region, destinationType).Set(c)
	}
	for sendingPool, c := range configSetByPoolCounts {
		metrics.SESConfigSetsBySendingPool.WithLabelValues(account, accountName, region, sendingPool).Set(c)
	}

	var contactListCount float64
	contactListPaginator := sesv2.NewListContactListsPaginator(client, &sesv2.ListContactListsInput{})
	for contactListPaginator.HasMorePages() {
		page, err := contactListPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		contactListCount += float64(len(page.ContactLists))
	}
	metrics.SESContactLists.WithLabelValues(account, accountName, region).Set(contactListCount)

	suppressionCounts := make(map[string]float64)
	suppressedPaginator := sesv2.NewListSuppressedDestinationsPaginator(client, &sesv2.ListSuppressedDestinationsInput{})
	for suppressedPaginator.HasMorePages() {
		page, err := suppressedPaginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, destination := range page.SuppressedDestinationSummaries {
			reason := string(destination.Reason)
			if reason == "" {
				reason = "UNKNOWN"
			}
			suppressionCounts[reason]++
		}
	}
	for reason, c := range suppressionCounts {
		metrics.SESSuppressedDestinations.WithLabelValues(account, accountName, region, reason).Set(c)
	}

	dedicatedIPPoolCount := listDedicatedIPPools(ctx, client)
	metrics.SESDedicatedIPPools.WithLabelValues(account, accountName, region).Set(dedicatedIPPoolCount)

	accountOutput, err := client.GetAccount(ctx, &sesv2.GetAccountInput{})
	if err != nil {
		log.Warn().Err(err).Str("region", region).Msg("Failed to get SES account details")
		return nil
	}
	metrics.SESAccountSettings.WithLabelValues(account, accountName, region, "sending_enabled").Set(boolAsFloat(accountOutput.SendingEnabled))
	metrics.SESAccountSettings.WithLabelValues(account, accountName, region, "production_access_enabled").Set(boolAsFloat(accountOutput.ProductionAccessEnabled))
	metrics.SESAccountSettings.WithLabelValues(account, accountName, region, "dedicated_ip_auto_warmup_enabled").Set(boolAsFloat(accountOutput.DedicatedIpAutoWarmupEnabled))
	suppressionEnabled := false
	if accountOutput.SuppressionAttributes != nil {
		suppressionEnabled = len(accountOutput.SuppressionAttributes.SuppressedReasons) > 0
	}
	metrics.SESAccountSettings.WithLabelValues(account, accountName, region, "suppression_enabled").Set(boolAsFloat(suppressionEnabled))
	if accountOutput.SendQuota != nil {
		metrics.SESSendingQuota.WithLabelValues(account, accountName, region, "max_24_hour_send").Set(accountOutput.SendQuota.Max24HourSend)
		metrics.SESSendingQuota.WithLabelValues(account, accountName, region, "max_send_rate").Set(accountOutput.SendQuota.MaxSendRate)
		metrics.SESSendingQuota.WithLabelValues(account, accountName, region, "sent_last_24_hours").Set(accountOutput.SendQuota.SentLast24Hours)
	}
	return nil
}

func boolAsFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func eventDestinationTypes(destination sesv2Types.EventDestination) []string {
	types := make([]string, 0)
	if destination.CloudWatchDestination != nil {
		types = append(types, "cloudwatch")
	}
	if destination.EventBridgeDestination != nil {
		types = append(types, "eventbridge")
	}
	if destination.KinesisFirehoseDestination != nil {
		types = append(types, "firehose")
	}
	if destination.PinpointDestination != nil {
		types = append(types, "pinpoint")
	}
	if destination.SnsDestination != nil {
		types = append(types, "sns")
	}
	if len(types) == 0 {
		types = append(types, "unknown")
	}
	return types
}

func listDedicatedIPPools(ctx context.Context, client *sesv2.Client) float64 {
	var total float64
	paginator := sesv2.NewListDedicatedIpPoolsPaginator(client, &sesv2.ListDedicatedIpPoolsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return total
		}
		total += float64(len(page.DedicatedIpPools))
	}
	return total
}
