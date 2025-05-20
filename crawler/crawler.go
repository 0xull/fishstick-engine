package crawler

import (
	"context"
	"net/http"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"
	"github.com/IkehAkinyemi/fishstick-engine/pipeline"
)

//go:generate mockgen -package mocks -destination mocks/mocks.go github.com/IkehAkinyemi/fishstick-engine/crawler URLGetter,PrivateNetworkDetector,Graph,Indexer

// URLGetter is implemented by objects that can perform HTTP GET requests
type URLGetter interface {
	Get(url string) (*http.Response, error)
}

// PrivateNetworkDetector is implemented by objects that can detect whether a
// host resolves to a private network address.
type PrivateNetworkDetector interface {
	IsPrivate(host string) (bool, error)
}

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
