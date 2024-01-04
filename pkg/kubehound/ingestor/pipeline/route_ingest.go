package pipeline

import (
	"context"
	"github.com/DataDog/KubeHound/pkg/globals/types"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph/vertex"
	"github.com/DataDog/KubeHound/pkg/kubehound/ingestor/preflight"
	"github.com/DataDog/KubeHound/pkg/kubehound/storage/cache"
	"github.com/DataDog/KubeHound/pkg/kubehound/store/collections"
)

const (
	RouteIngestName = "openshift-route-ingest"
)

type RouteIngest struct {
	vertex     *vertex.Route
	collection collections.Route
	r          *IngestResources
}

func (i *RouteIngest) Name() string { return RouteIngestName }

var _ ObjectIngest = (*NodeIngest)(nil)

func (i *RouteIngest) Initialize(ctx context.Context, deps *Dependencies) error {
	var err error

	i.vertex = &vertex.Route{}
	i.collection = collections.Route{}

	i.r, err = CreateResources(ctx, deps,
		WithCacheWriter(cache.WithExpectedOverwrite()),
		WithStoreWriter(i.collection),
		WithGraphWriter(i.vertex))
	if err != nil {
		return err
	}

	return nil
}

func (i *RouteIngest) IngestRoute(ctx context.Context, r types.RouteType) error {
	if ok, err := preflight.CheckRoute(r); !ok {
		return err
	}

	// Normalize node to store object format
	o, err := i.r.storeConvert.Route(ctx, r)
	if err != nil {
		return err
	}

	// Async write to store
	if err := i.r.writeStore(ctx, i.collection, o); err != nil {
		return err
	}

	// Transform store model to vertex input
	insert, err := i.r.graphConvert.Route(o)
	if err != nil {
		return err
	}

	return i.r.writeVertex(ctx, i.vertex, insert)
}

// Complete is invoked by the collector when all nodes have been streamed.
// The function flushes all writers and waits for completion.
func (i *RouteIngest) Complete(ctx context.Context) error {
	return i.r.flushWriters(ctx)
}

func (i *RouteIngest) Run(ctx context.Context) error {
	return i.r.collect.StreamRoutes(ctx, i)
}

func (i *RouteIngest) Close(ctx context.Context) error {
	return i.r.cleanupAll(ctx)
}
