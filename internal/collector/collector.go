package collector

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/nimishgj/aws-radar/internal/metrics"
	"github.com/rs/zerolog/log"
)

// Collector interface for all AWS service collectors
type Collector interface {
	Name() string
	Collect(ctx context.Context, cfg aws.Config, region string) error
}

// GlobalCollector interface for services that don't use regions (IAM, Route53, etc.)
type GlobalCollector interface {
	Name() string
	Collect(ctx context.Context, cfg aws.Config) error
}

// Orchestrator manages all collectors and runs collection cycles
type Orchestrator struct {
	collectors       []Collector
	globalCollectors []GlobalCollector
	regions          []string
	interval         time.Duration
	timeout          time.Duration
}

// NewOrchestrator creates a new collector orchestrator
func NewOrchestrator(regions []string, interval, timeout time.Duration, enabledCollectors []string) *Orchestrator {
	allCollectors := []Collector{
		NewAPIGatewayCollector(),
		NewAPIGatewayV2Collector(),
		NewAutoScalingCollector(),
		NewAthenaCollector(),
		NewECRCollector(),
		NewEC2Collector(),
		NewEFSCollector(),
		NewEventBridgeCollector(),
		NewGlueCollector(),
		NewS3Collector(),
		NewRDSCollector(),
		NewLambdaCollector(),
		NewECSCollector(),
		NewELBCollector(),
		NewEKSCollector(),
		NewDynamoDBCollector(),
		NewElastiCacheCollector(),
		NewOpenSearchCollector(),
		NewSecretsManagerCollector(),
		NewSfnCollector(),
		NewSSMCollector(),
		NewSQSCollector(),
		NewSNSCollector(),
		NewEBSCollector(),
		NewVPCCollector(),
		NewACMCollector(),
	}
	allGlobalCollectors := []GlobalCollector{
		NewCloudFrontCollector(),
		NewRoute53Collector(),
		NewIAMCollector(),
	}

	enabled := normalizeCollectors(enabledCollectors)
	if len(enabled) > 0 {
		known := make(map[string]struct{}, len(allCollectors)+len(allGlobalCollectors))
		for _, c := range allCollectors {
			known[c.Name()] = struct{}{}
		}
		for _, c := range allGlobalCollectors {
			known[c.Name()] = struct{}{}
		}
		for name := range enabled {
			if _, ok := known[name]; !ok {
				log.Warn().Str("collector", name).Msg("Unknown collector configured")
			}
		}
	}

	collectors := filterCollectors(allCollectors, enabled)
	globalCollectors := filterGlobalCollectors(allGlobalCollectors, enabled)

	if len(collectors) == 0 && len(globalCollectors) == 0 {
		log.Warn().Msg("No collectors enabled; nothing will be collected")
	}

	return &Orchestrator{
		collectors:       collectors,
		globalCollectors: globalCollectors,
		regions:          regions,
		interval:         interval,
		timeout:          timeout,
	}
}

func normalizeCollectors(enabled []string) map[string]struct{} {
	normalized := make(map[string]struct{})
	for _, name := range enabled {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		normalized[name] = struct{}{}
	}
	return normalized
}

func filterCollectors(all []Collector, enabled map[string]struct{}) []Collector {
	if len(enabled) == 0 {
		return all
	}
	filtered := make([]Collector, 0, len(all))
	for _, c := range all {
		if _, ok := enabled[c.Name()]; ok {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func filterGlobalCollectors(all []GlobalCollector, enabled map[string]struct{}) []GlobalCollector {
	if len(enabled) == 0 {
		return all
	}
	filtered := make([]GlobalCollector, 0, len(all))
	for _, c := range all {
		if _, ok := enabled[c.Name()]; ok {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// Start begins the collection loop
func (o *Orchestrator) Start(ctx context.Context) {
	log.Info().
		Strs("regions", o.regions).
		Dur("interval", o.interval).
		Msg("Starting collector orchestrator")

	// Run initial collection
	o.collect(ctx)

	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping collector orchestrator")
			return
		case <-ticker.C:
			o.collect(ctx)
		}
	}
}

func (o *Orchestrator) collect(ctx context.Context) {
	log.Info().Msg("Starting collection cycle")
	start := time.Now()

	// Reset all metrics before collecting
	metrics.ResetAll()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS config")
		return
	}

	var wg sync.WaitGroup

	// Run regional collectors
	for _, region := range o.regions {
		for _, collector := range o.collectors {
			wg.Add(1)
			go func(c Collector, r string) {
				defer wg.Done()
				o.runCollector(ctx, c, cfg, r)
			}(collector, region)
		}
	}

	// Run global collectors
	for _, collector := range o.globalCollectors {
		wg.Add(1)
		go func(c GlobalCollector) {
			defer wg.Done()
			o.runGlobalCollector(ctx, c, cfg)
		}(collector)
	}

	wg.Wait()

	duration := time.Since(start)
	log.Info().
		Dur("duration", duration).
		Msg("Collection cycle completed")
}

func (o *Orchestrator) runCollector(ctx context.Context, c Collector, cfg aws.Config, region string) {
	collectorCtx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	start := time.Now()
	regionCfg := cfg.Copy()
	regionCfg.Region = region

	if err := c.Collect(collectorCtx, regionCfg, region); err != nil {
		log.Error().
			Err(err).
			Str("collector", c.Name()).
			Str("region", region).
			Msg("Collection failed")
		metrics.CollectionUp.WithLabelValues(c.Name(), region).Set(0)
		metrics.CollectionErrors.WithLabelValues(c.Name(), region).Inc()
	} else {
		metrics.CollectionUp.WithLabelValues(c.Name(), region).Set(1)
	}

	duration := time.Since(start).Seconds()
	metrics.CollectionDuration.WithLabelValues(c.Name()).Observe(duration)
}

func (o *Orchestrator) runGlobalCollector(ctx context.Context, c GlobalCollector, cfg aws.Config) {
	collectorCtx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	start := time.Now()

	if err := c.Collect(collectorCtx, cfg); err != nil {
		log.Error().
			Err(err).
			Str("collector", c.Name()).
			Msg("Global collection failed")
		metrics.CollectionUp.WithLabelValues(c.Name(), "global").Set(0)
		metrics.CollectionErrors.WithLabelValues(c.Name(), "global").Inc()
	} else {
		metrics.CollectionUp.WithLabelValues(c.Name(), "global").Set(1)
	}

	duration := time.Since(start).Seconds()
	metrics.CollectionDuration.WithLabelValues(c.Name()).Observe(duration)
}
