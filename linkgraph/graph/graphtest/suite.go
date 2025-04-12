package graphtest

import (
	"time"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"
	"github.com/google/uuid"
	gc "gopkg.in/check.v1"
)

// SuiteBase defines a re-usable set of graph-related tests that can
// be executed against any type that implements graph.Graph
type SuiteBase struct {
	g graph.Graph
}

// SetGraph configures the test-suite to run all tests against g.
func (s *SuiteBase) SetGraph(g graph.Graph) {
	s.g = g
}

// TestUpsertLink verifies the link upsert logic
func (s *SuiteBase) TestUpsertLink(c *gc.C) {
	// Create a new link
	original := &graph.Link{
		URL: "https://example.com",
		RetrievedAt: time.Now().Add(-10 * time.Hour),
	}
	
	err := s.g.UpsertLink(original)
	c.Assert(err, gc.IsNil)
	c.Assert(original.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf("expected a linkID to be assigned to the new link"))
	
	// Update existing link with a newer timestamp and different URL
	accessedAt := time.Now().Truncate(time.Second).UTC()
	existing := &graph.Link{
		ID: original.ID,
		URL: "https://example.com",
		RetrievedAt: accessedAt,
	}
	err = s.g.UpsertLink(existing)
	c.Assert(err, gc.IsNil)
	c.Assert(existing.ID, gc.Equals, original.ID, gc.Commentf("link ID changed while upserting"))
}
