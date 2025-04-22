package inmemory

import "github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"

// LinkIterator is a graph.LinkIterator implementation for the in-memory graph.
type linkIterator struct {
	s *InMemoryGraph
	
	links []*graph.Link
	currIndex int
}

// Next implements graph.LinkIterator.
func (i *linkIterator) Next() bool {
	if i.currIndex >= len(i.links) {
		return false
	}
	i.currIndex++
	return true
}

// Error implements graph.LinkIterator
func (i *linkIterator) Error() error {
	return nil
}

// Close implements graph.LinkIterator.
func (i *linkIterator) Close() error {
	return nil
}

// Link implements graph.LinkIterator.
func (i linkIterator) Link() *graph.Link {
	// The link pointer contents may be overwritten by graph update;
	// to avoid data-races we acquire the read lock first and clone the Link
	i.s.mu.RLock()
	link := new(graph.Link)
	*link = *i.links[i.currIndex-1]
	i.s.mu.RUnlock()
	return link
}

// edgeIterator is a graph.EdgeIterator implementation for the in-memory graph.
type edgeIterator struct {
	s *InMemoryGraph
	
	edges []*graph.Edge
	currIndex int
}

// Next implements graph.EdgeIterator.
func (e edgeIterator) Next() bool {
	if i.currIndex >= len(e.edges) {
		return false
	}
	i.currIndex++
	return true
}

// Error implements graph.EdgeIterator.
func (e edgeIterator) Error() error {
	return nil
}

// Close implements graph.EdgeIterator.
func (e edgeIterator) Close() error {
	return nil
}

// Link implements graph.EdgeIterator.
func (e edgeIterator) Edge() *graph.Edge {
	// The edge pointer contents may be overwritten by a graph update; to
	// avoid data-races we acquire the read lock first and clone the Edge
	e.s.mu.RLock()
	edge := new(graph.Edge)
	*edge = *e.edges[e.currIndex-1]
	e.s.mu.RUnlock()
	return edge
}
