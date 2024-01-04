package vertex

import (
	"fmt"
	"github.com/DataDog/KubeHound/pkg/kubehound/graph/types"
	"github.com/DataDog/KubeHound/pkg/kubehound/models/graph"
	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_Traversal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want types.VertexTraversal
		data graph.Route
	}{
		{
			name: "Add Identities in JanusGraph",
			// We set the values to all field with non default values
			// so we are sure all are correctly propagated.
			data: graph.Route{
				StoreID:      "TestStoreID",
				App:          "TestApp",
				Team:         "TestTeam",
				Service:      "TestService",
				RunID:        "TestRunID",
				Cluster:      "TestCluster",
				IsNamespaced: true,
				Namespace:    "TestNamespace",
				Name:         "TestName",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Route{}

			g := gremlingo.GraphTraversalSource{}

			routeTraversal := r.Traversal()
			inserts := []any{&tt.data}

			traversal := routeTraversal(&g, inserts)
			// This is ugly but doesn't need to write to the DB
			// This just makes sure the traversal is correctly returned with the correct values
			assert.Contains(t, fmt.Sprintf("%s", traversal.Bytecode), "TestStoreID")
			assert.Contains(t, fmt.Sprintf("%s", traversal.Bytecode), "TestNamespace")
			assert.Contains(t, fmt.Sprintf("%s", traversal.Bytecode), "TestName")
			assert.Contains(t, fmt.Sprintf("%s", traversal.Bytecode), "TestService")
		})
	}
}
