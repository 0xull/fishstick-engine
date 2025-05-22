package crawler

import (
	"context"
	"net/http"
	"time"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"
	"github.com/IkehAkinyemi/fishstick-engine/pipeline"
	"github.com/IkehAkinyemi/fishstick-engine/textindexer/index"
	"github.com/google/uuid"
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

// Graph is implemented by objects that can upsert links and edges into a link
// graph instance.
type Graph interface {
	// UpsertLink creates a new link or updates an existing link.
	UpsertLink(link *graph.Link) error
	
	// UpsertEdge creates a new edge or updates an existing edge.
	UpsertEdge(edge *graph.Edge) error
	
	// RemoveStaleEdges removes any edge that originates from the specified
	// link ID and was updated before the specified timestamp.
	RemoveStaleEdges(fromID uuid.UUID, updatedBefore time.Time) error
}

// Indexer is implemented by objects that can index the contents of webpages
// retrieved by the crawler pipeline.
type Indexer interface {
	// Index inserts a new document to the index or updates the index entry
	// for and existing document.
	Index(doc *index.Document) error
}

// Config encapsulates the configuration options for creating a new Crawler.
type Config struct {
	// A PrivateNetworkDetector instance
	PrivateNetworkDetector PrivateNetworkDetector
	
	// A URLGetter instance for fetching links.
	URLGetter URLGetter
	
	// A GraphUpdater instance for adding new links to the link graph.
	Graph Graph
	
	// A TextIndexer instance for indexing the content of each retrieved link.
	Indexer Indexer
	
	// The number of concurrent worker used for retrieving links.
	FetchWorkers int
}

// Crawler implements a webpage crawling pipeline consisting of the following
// stage:
// 	 - Given a URL, retrieve the webpage contents from the remote server
//   - Extract and resolve absolute and relative links from the retrieved page.
//   - Extract page title and text content from the retrieved page.
//   - Update the link graph: add new links and create edges between crawled
// 	   page and the links within it.
//   - Index crawled page title and text content
type Crawler struct {
	p *pipeline.Pipeline
}
