package core

import (
	"context"
	"fmt"
	"time"

	"github.com/DataDog/KubeHound/pkg/collector"
	"github.com/DataDog/KubeHound/pkg/config"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph/edge"
	"github.com/DataDog/KubeHound/pkg/kubehound/ingestor"
	"github.com/DataDog/KubeHound/pkg/kubehound/storage/cache"
	"github.com/DataDog/KubeHound/pkg/kubehound/storage/graphdb"
	"github.com/DataDog/KubeHound/pkg/kubehound/storage/storedb"
	"github.com/DataDog/KubeHound/pkg/telemetry"
	"github.com/DataDog/KubeHound/pkg/telemetry/log"
	"github.com/DataDog/KubeHound/pkg/telemetry/span"
	"github.com/DataDog/KubeHound/pkg/telemetry/tag"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func ingestData(ctx context.Context, cfg *config.KubehoundConfig, collect collector.CollectorClient, cache cache.CacheProvider,
	storedb storedb.Provider, graphdb graphdb.Provider) error {

	start := time.Now()
	span, ctx := tracer.StartSpanFromContext(ctx, span.IngestData, tracer.Measured())
	defer span.Finish()

	log.I.Info("Loading data ingestor")
	ingest, err := ingestor.Factory(cfg, collect, cache, storedb, graphdb)
	if err != nil {
		return fmt.Errorf("ingestor creation: %w", err)
	}
	defer ingest.Close(ctx)

	log.I.Info("Running dependency health checks")
	if err := ingest.HealthCheck(ctx); err != nil {
		return fmt.Errorf("ingestor dependency health check: %w", err)
	}

	log.I.Info("Running data ingest and normalization")
	if err := ingest.Run(ctx); err != nil {
		return fmt.Errorf("ingest: %w", err)
	}

	log.I.Infof("Completed data ingest and normalization in %s", time.Since(start))

	return nil
}

// buildGraph will construct the attack graph by calculating and inserting all registered edges in parallel.
// All I/O operations are performed asynchronously.
func buildGraph(ctx context.Context, cfg *config.KubehoundConfig, storedb storedb.Provider,
	graphdb graphdb.Provider, cache cache.CacheReader) error {

	start := time.Now()
	span, ctx := tracer.StartSpanFromContext(ctx, span.BuildGraph, tracer.Measured())
	defer span.Finish()

	log.I.Info("Loading graph edge definitions")
	edges := edge.Registered()
	if err := edges.Verify(); err != nil {
		return fmt.Errorf("edge registry verification: %w", err)
	}

	log.I.Info("Loading graph builder")
	builder, err := graph.NewBuilder(cfg, storedb, graphdb, cache, edges)
	if err != nil {
		return fmt.Errorf("graph builder creation: %w", err)
	}

	log.I.Info("Running dependency health checks")
	if err := builder.HealthCheck(ctx); err != nil {
		return fmt.Errorf("graph builder dependency health check: %w", err)
	}

	log.I.Info("Constructing graph")
	if err := builder.Run(ctx); err != nil {
		return fmt.Errorf("graph builder edge calculation: %w", err)
	}

	log.I.Infof("Completed graph construction in %s", time.Since(start))

	return nil
}

// Launch will launch the KubeHound application to ingest data from a collector and create an attack graph.
func Launch(ctx context.Context, opts ...LaunchOption) error {
	span, ctx := tracer.StartSpanFromContext(ctx, span.Launch, tracer.Measured())
	defer span.Finish()

	// We define a unique run id this so we can measure run by run in addition of version per version.
	// Useful when rerunning the same binary (same version) on different dataset or with different databases...
	runID := config.NewRunID()
	span.SetBaggageItem("run_id", runID.String())

	// We update the base tags to include that run id, so we have it available for metrics
	tag.BaseTags = append(tag.BaseTags, tag.RunID(runID.String()))

	// Set the run ID as a global log tag
	log.AddGlobalTags(map[string]string{
		tag.RunIdTag: runID.String(),
	})

	// Start the run
	start := time.Now()
	log.I.Infof("Starting KubeHound (run_id: %s)", runID.String())
	log.I.Info("Initializing launch options")
	lOpts := &launchConfig{}
	for _, opt := range opts {
		opt(lOpts)
	}

	var cfg *config.KubehoundConfig
	if len(lOpts.ConfigPath) != 0 {
		log.I.Infof("Loading application configuration from file %s", lOpts.ConfigPath)
		cfg = config.MustLoadConfig(lOpts.ConfigPath)
	} else {
		log.I.Infof("Loading application configuration from default embedded")
		cfg = config.MustLoadEmbedConfig()
	}

	// Update the logger behaviour from configuration
	log.SetDD(cfg.Telemetry.Enabled)
	log.AddGlobalTags(cfg.Telemetry.Tags)

	// Setup telemetry
	log.I.Info("Initializing application telemetry")
	ts, err := telemetry.Initialize(cfg)
	if err != nil {
		log.I.Warnf("failed telemetry initialization: %v", err)
	}
	defer telemetry.Shutdown(ts)

	// Create the cache client
	log.I.Info("Loading cache provider")
	cp, err := cache.Factory(ctx, cfg)
	if err != nil {
		return fmt.Errorf("cache client creation: %w", err)
	}
	defer cp.Close(ctx)
	log.I.Infof("Loaded %s cache provider", cp.Name())

	// Create the store client
	log.I.Info("Loading store database provider")
	sp, err := storedb.Factory(ctx, cfg)
	if err != nil {
		return fmt.Errorf("store database client creation: %w", err)
	}
	defer sp.Close(ctx)
	log.I.Infof("Loaded %s store provider", sp.Name())

	err = sp.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("store database prepare: %w", err)
	}

	// Create the graph client
	log.I.Info("Loading graph database provider")
	gp, err := graphdb.Factory(ctx, cfg)
	if err != nil {
		return fmt.Errorf("graph database client creation: %w", err)
	}
	defer gp.Close(ctx)
	log.I.Infof("Loaded %s graph provider", gp.Name())

	err = gp.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("graph database prepare: %w", err)
	}

	// Create the collector instance
	log.I.Info("Loading Kubernetes data collector client")
	collect, err := collector.ClientFactory(ctx, cfg)
	if err != nil {
		return fmt.Errorf("collector client creation: %w", err)
	}
	defer collect.Close(ctx)
	log.I.Infof("Loaded %s collector client", collect.Name())

	// All dependencies are loaded - we can complete the config with runtime information
	cluster, err := collect.ClusterInfo(ctx)
	if err != nil {
		return fmt.Errorf("collector cluster info: %w", err)
	}
	cfg.ComputeDynamic(config.WithClusterName(cluster.Name))

	// Run the ingest pipeline
	log.I.Info("Starting Kubernetes raw data ingest")
	if err := ingestData(ctx, cfg, collect, cp, sp, gp); err != nil {
		return fmt.Errorf("raw data ingest: %w", err)
	}

	// Construct the graph
	log.I.Info("Building attack graph")
	if err := buildGraph(ctx, cfg, sp, gp, cp); err != nil {
		return fmt.Errorf("building attack graph: %w", err)
	}

	log.I.Infof("KubeHound run (id=%s) complete in %s", runID.String(), time.Since(start))

	return nil
}
