package collector

import (
	"context"
	"fmt"
	"github.com/DataDog/KubeHound/pkg/config"
	"github.com/DataDog/KubeHound/pkg/globals"
	"github.com/DataDog/KubeHound/pkg/telemetry/log"
	"github.com/DataDog/KubeHound/pkg/telemetry/metric"
	"github.com/DataDog/KubeHound/pkg/telemetry/span"
	"github.com/DataDog/KubeHound/pkg/telemetry/statsd"
	"github.com/DataDog/KubeHound/pkg/telemetry/tag"
	routev1 "github.com/openshift/api/route/v1"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io/fs"
	"path/filepath"
)

const (
	routePath = "routes.route.openshift.io.json"
)

const (
	FileOpenshiftCollectorName = "local-file-openshift-collector"
)

type openShiftFileCollector struct {
	*FileCollector
}

func NewOpenShiftFileCollector(ctx context.Context, cfg *config.KubehoundConfig) (CollectorClient, error) {
	tags := tag.BaseTags
	tags = append(tags, tag.Collector(OpenShiftAPICollectorName))
	if cfg.Collector.Type != config.CollectorTypeOpenShiftFile {
		return nil, fmt.Errorf("invalid collector type in config: %s", cfg.Collector.Type)
	}

	l := log.Trace(ctx, log.WithComponent(globals.FileCollectorComponent))
	l.Infof("Creating file collector from directory %s", cfg.Collector.File.Directory)

	return &openShiftFileCollector{
		FileCollector: &FileCollector{
			cfg:  cfg.Collector.File,
			log:  l,
			tags: tags,
		},
	}, nil
}

func (c *openShiftFileCollector) Name() string {
	return FileOpenshiftCollectorName
}

func (c *openShiftFileCollector) streamRoutesNamespace(ctx context.Context, fp string, ingestor RouteIngestor) error {
	list, err := readList[routev1.RouteList](ctx, fp)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		_ = statsd.Incr(metric.CollectorCount, append(c.tags, tag.Entity(tag.EntityRoutes)), 1)
		i := item
		err = ingestor.IngestRoute(ctx, &i)
		if err != nil {
			return fmt.Errorf("processing OpenShift route %s: %w", i.Name, err)
		}
	}

	return nil
}

func (c *openShiftFileCollector) StreamRoutes(ctx context.Context, ingestor RouteIngestor) error {
	span, ctx := tracer.StartSpanFromContext(ctx, span.CollectorStream, tracer.Measured())
	span.SetTag(tag.EntityTag, tag.EntityRoutes)
	defer span.Finish()

	err := filepath.WalkDir(c.cfg.Directory, func(path string, d fs.DirEntry, err error) error {
		if path == c.cfg.Directory || !d.IsDir() {
			// Skip files
			return nil
		}

		fp := filepath.Join(path, routePath)
		c.log.Debugf("Streaming pods from file %s", fp)

		return c.streamRoutesNamespace(ctx, fp, ingestor)
	})

	if err != nil {
		return fmt.Errorf("file collector stream routes: %w", err)
	}

	return ingestor.Complete(ctx)
}
