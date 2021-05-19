package quadgraph_test

import (
	"fmt"
	"quadgraph"
	"testing"

	"github.com/cayleygraph/quad"
)

func TestStore(t *testing.T) {
	qs := quadgraph.NewStore()

	qs.AddQuad(quad.Make(quad.IRI("oxisto"), "creates", "database", nil))
	qs.AddQuad(quad.Make(quad.IRI("oxisto"), "cannot", "sleep", nil))

	fmt.Printf("%+v", qs)
}

func TestPeformance(t *testing.T) {
	qs := quadgraph.NewStore()

	for i := 0; i < 1000000; i++ {
		qs.AddQuad(quad.Make(quad.IRI("oxisto"), fmt.Sprintf("creates%d", i), "database", nil))
	}

	fmt.Printf("%+v", qs)
}
