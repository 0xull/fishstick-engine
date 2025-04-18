package graphtest

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
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
		URL:         "https://example.com",
		RetrievedAt: time.Now().Add(-10 * time.Hour),
	}

	err := s.g.UpsertLink(original)
	c.Assert(err, gc.IsNil)
	c.Assert(original.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf("expected a linkID to be assigned to the new link"))

	// Update existing link with a newer timestamp and different URL
	accessedAt := time.Now().Truncate(time.Second).UTC()
	existing := &graph.Link{
		ID:          original.ID,
		URL:         "https://example.com",
		RetrievedAt: accessedAt,
	}
	err = s.g.UpsertLink(existing)
	c.Assert(err, gc.IsNil)
	c.Assert(existing.ID, gc.Equals, original.ID, gc.Commentf("link ID changed while upserting"))

	stored, err := s.g.FindLink(existing.ID)
	c.Assert(err, gc.IsNil)
	c.Assert(stored.RetrievedAt, gc.Equals, accessedAt, gc.Commentf("last accessed timestamp was not updated"))

	// Attempt to insert a new link whose URL matches an existng link with
	// and provide an older accessedAt value
	sameURL := &graph.Link{
		URL:         existing.URL,
		RetrievedAt: time.Now().Add(-10 * time.Hour).UTC(),
	}

	err = s.g.UpsertLink(sameURL)
	c.Assert(err, gc.IsNil)
	c.Assert(sameURL.ID, gc.Equals, existing.ID)

	stored, err = s.g.FindLink(existing.ID)
	c.Assert(err, gc.IsNil)
	c.Assert(stored.RetrievedAt, gc.Equals, accessedAt, gc.Commentf("last accessed timestamp was overwritten with older value"))

	// Create a new link and then attempt to update its URL to the same as
	// an existing link.
	dup := &graph.Link{
		URL: "foo",
	}
	err = s.g.UpsertLink(dup)
	c.Assert(err, gc.IsNil)
	c.Assert(dup.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf("expected a linkID to be assigned to the new link"))
}

// TestFindLink verifies the link lookup logic.
func (s *SuiteBase) TestFindLink(c *gc.C) {
	// Create a new link
	link := &graph.Link{
		URL:         "https://example.com",
		RetrievedAt: time.Now().Truncate(time.Second).UTC(),
	}

	err := s.g.UpsertLink(link)
	c.Assert(err, gc.IsNil)
	c.Assert(link.ID, gc.Not(gc.Equals), uuid.Nil, gc.Commentf("expected a linkID to be assigned to the new link"))

	// lookup link by ID
	other, err := s.g.FindLink(link.ID)
	c.Assert(err, gc.IsNil)
	c.Assert(other, gc.DeepEquals, link, gc.Commentf("lookup by ID returned the wrong link"))

	// lookup link by unknown ID
	_, err = s.g.FindLink(uuid.Nil)
	c.Assert(errors.Is(err, graph.ErrNotFound), gc.Equals, true)
}

// TestConcurrentLinkIterators verifies that multiple clients can concurrently
// access the store.
func (s *SuiteBase) TestConcurrentLinkIterators(c *gc.C) {
	var (
		wg           sync.WaitGroup
		numIterators = 10
		numLinks     = 100
	)

	for i := range numLinks {
		link := &graph.Link{URL: fmt.Sprint(i)}
		c.Assert(s.g.UpsertLink(link), gc.IsNil)
	}

	wg.Add(numIterators)
	for i := range numIterators {
		go func(id int) {
			defer wg.Done()

			itTagComment := gc.Commentf("iterator %d", id)
			seen := make(map[string]bool)
			it, err := s.partitionedLinkIterator(c, 0, 1, time.Now())
			c.Assert(err, gc.IsNil, itTagComment)
			defer func() {
				c.Assert(it.Close(), gc.IsNil, itTagComment)
			}()

			for it.Next() {
				link := it.Link()
				linkID := link.ID.String()
				c.Assert(seen[linkID], gc.Equals, false, gc.Commentf("iterator %d saw same link twice", id))
			}

			c.Assert(seen, gc.HasLen, numLinks, itTagComment)
			c.Assert(it.Error(), gc.IsNil, itTagComment)
			c.Assert(it.Close(), gc.IsNil, itTagComment)
		}(i)
	}

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		//test completed successfully
	case <-time.After(10 * time.Second):
		c.Fatalf("timed out waiting for test to complete")
	}
}

func (s *SuiteBase) partitionedLinkIterator(c *gc.C, partition, numPartitions int, accessedBefore time.Time) (graph.LinkIterator, error) {
	from, to := s.partitionRange(c, partition, numPartitions)
	return s.g.Links(from, to, accessedBefore)
}

func (s *SuiteBase) partitionRange(c *gc.C, partition, numPartitions int) (from, to uuid.UUID) {
	if partition < 0 || partition >= numPartitions {
		c.Fatal("invalid partition")
	}

	var minUUID = uuid.Nil
	var maxUUID = uuid.MustParse("ffffffff-ffff-ffff-ffffffffffff")
	var err error

	// Calculate the size of each partition as: (2^128 / numPartitions)
	tokenRange := big.NewInt(0)
	partSize := big.NewInt(0)
	partSize.SetBytes(maxUUID[:])
	partSize = partSize.Div(partSize, big.NewInt(int64(numPartitions)))

	// We model the partitions as a segment that begins at minUUID (all
	// bits set to zero) and ends maxUUID (all bits set to 1). By
	// setting the end range for the *last* partition to maxUUID we ensure
	// that we always cover the full range of UUIDs even if the range
	// itself is not evenly divisible by numPartitions.
	if partition == 0 {
		from = minUUID
	} else {
		tokenRange.Mul(partSize, big.NewInt(int64(partition)))
		from, err = uuid.FromBytes(tokenRange.Bytes())
		c.Assert(err, gc.IsNil)
	}

	if partition == numPartitions-1 {
		to = maxUUID
	} else {
		tokenRange.Mul(partSize, big.NewInt(int64(partition+1)))
		to, err = uuid.FromBytes(tokenRange.Bytes())
		c.Assert(err, gc.IsNil)
	}

	return from, to
}
