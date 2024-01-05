package pipeline

import (
	"context"
	"github.com/DataDog/KubeHound/pkg/collector"
	mockcollect "github.com/DataDog/KubeHound/pkg/collector/mockcollector"
	"github.com/DataDog/KubeHound/pkg/config"
	"github.com/DataDog/KubeHound/pkg/globals/types"
	"github.com/DataDog/KubeHound/pkg/kubehound/models/store"
	cache "github.com/DataDog/KubeHound/pkg/kubehound/storage/cache/mocks"
	storedb "github.com/DataDog/KubeHound/pkg/kubehound/storage/storedb/mocks"
	"github.com/DataDog/KubeHound/pkg/kubehound/store/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRouteIngest_Pipeline(t *testing.T) {
	t.Parallel()

	ri := &RouteIngest{}
	ctx := context.Background()
	fakeRoute, err := loadTestObject[types.RouteType]("testdata/route.json")
	assert.NoError(t, err)

	client := mockcollect.NewOpenShiftCollectorClient(t)
	client.EXPECT().StreamRoutes(ctx, ri).
		RunAndReturn(func(ctx context.Context, i collector.RouteIngestor) error {
			// Fake the stream of a single route from the collector client
			err := i.IngestRoute(ctx, fakeRoute)
			if err != nil {
				return err
			}

			return i.Complete(ctx)
		})

	// Cache setup
	c := cache.NewCacheProvider(t)
	cw := cache.NewAsyncWriter(t)

	cw.EXPECT().Queue(ctx, mock.AnythingOfType("*cachekey.routeCacheKey"), mock.AnythingOfType("store.Route")).Return(nil).Once()
	cw.EXPECT().Flush(ctx).Return(nil)
	cw.EXPECT().Close(ctx).Return(nil)
	c.EXPECT().BulkWriter(ctx).Return(cw, nil)

	// Store setup
	sdb := storedb.NewProvider(t)
	sw := storedb.NewAsyncWriter(t)
	route := collections.Route{}
	storeID := store.ObjectID()
	sw.EXPECT().Queue(ctx, mock.AnythingOfType("*store.Route")).
		RunAndReturn(func(ctx context.Context, i any) error {
			i.(*store.Route).Id = storeID

			return nil
		}).Once()
	sw.EXPECT().Flush(ctx).Return(nil)
	sw.EXPECT().Close(ctx).Return(nil)
	sdb.EXPECT().BulkWriter(ctx, route, mock.Anything).Return(sw, nil)

	deps := &Dependencies{
		Collector: client,
		Cache:     c,
		StoreDB:   sdb,
		Config: &config.KubehoundConfig{
			Builder: config.BuilderConfig{
				Edge: config.EdgeBuilderConfig{},
			},
			Dynamic: config.DynamicConfig{
				RunID:   testID,
				Cluster: "test-cluster",
			},
		},
	}

	// Initialize
	err = ri.Initialize(ctx, deps)
	assert.NoError(t, err)

	// Run
	err = ri.Run(ctx)
	assert.NoError(t, err)

	// Close
	err = ri.Close(ctx)
	assert.NoError(t, err)
}
