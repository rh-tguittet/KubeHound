package vertex

import (
	"context"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph/adapter"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph/types"
	"github.com/DataDog/KubeHound/pkg/kubehound/models/graph"
)

const (
	RouteLabel = "Route"
)

var _ Builder = (*Route)(nil)

type Route struct {
	BaseVertex
}

func (v *Route) Label() string { return RouteLabel }

func (v *Route) Processor(ctx context.Context, entry any) (any, error) {
	return adapter.GremlinVertexProcessor[*graph.Route](ctx, entry)
}

func (v *Route) Traversal() types.VertexTraversal {
	return v.DefaultTraversal(v.Label())
}
