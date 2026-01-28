package collector

import (
	"context"
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
func NewOrchestrator(regions []string, interval, timeout time.Duration) *Orchestrator {
	return &Orchestrator{
		collectors: []Collector{
			NewEC2Collector(),
			NewS3Collector(),
			NewRDSCollector(),
			NewLambdaCollector(),
			NewECSCollector(),
			NewELBCollector(),
			NewEKSCollector(),
			NewDynamoDBCollector(),
			NewElastiCacheCollector(),
			NewSQSCollector(),
			NewSNSCollector(),
			NewEBSCollector(),
			NewVPCCollector(),
			NewACMCollector(),
		},
		globalCollectors: []GlobalCollector{
			NewCloudFrontCollector(),
			NewRoute53Collector(),
			NewIAMCollector(),
		},
		regions:  regions,
		interval: interval,
		timeout:  timeout,
	}
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
		metrics.CollectionErrors.WithLabelValues(c.Name(), region).Inc()
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
		metrics.CollectionErrors.WithLabelValues(c.Name(), "global").Inc()
	}

	duration := time.Since(start).Seconds()
	metrics.CollectionDuration.WithLabelValues(c.Name()).Observe(duration)
}
