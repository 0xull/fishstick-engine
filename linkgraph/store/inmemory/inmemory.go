package inmemory

import (
	"sync"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph"
	"github.com/google/uuid"
)

// TODO: check InMemoryGraph implementation
// matches the Graph interface contract

type edgeList []uuid.UUID

type InMemoryGraph struct {
	mu sync.RWMutex
	
	links map[uuid.UUID]*graph.Link
	
}
