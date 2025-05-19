package crawler

import (
	"context"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"
	"github.com/IkehAkinyemi/fishstick-engine/pipeline"
)

type linkSource struct {
	linkIt graph.LinkIterator
}

func (ls *linkSource) Error() error { return ls.linkIt.Error() }
func (ls *linkSource) Next(context.Context) bool { return ls.linkIt.Next() }
func (ls *linkSource) Payload() pipeline.Payload {
	link := ls.linkIt.Link()
	p := payloadPool.Get().(*crawlerPayload)
	
	p.LinkID = link.ID
	p.URL = link.URL
	p.RetrievedAt = link.RetrievedAt
	return p
}

type countingSink struct {
	count int
}

func (s *countingSink) Consume(_ context.Context, p pipeline.Payload) error {
	s.count++
	return nil
}

func (s *countingSink) getCount() int {
	// The broadcast split-stage sends out two payloads for each incoming links
	// so we need to divide the total count by 2.
	return s.count / 2
}
