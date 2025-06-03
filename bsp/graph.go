package bspgraph

import "github.com/IkehAkinyemi/fishstick-engine/bsp/message"

// Vertex represents a vertex in the Graph.
type Vertex struct {
	id string
	value any
	active bool
	msgQueue [2]message.Queue
	edges []*Edge
}

// ID returns the vertex ID.
func (v *Vertex) ID() string { return v.id }

// Edges returns the list of outgoing edges from this vertex.
func (v *Vertex) Edges() []*Edge { return v.edges }

// Freeze marks the vertex as inactive. Inactive vertices will not be processed
// in the following supersteps unless they receive a message in which case they
// will be re-activated.
func (v *Vertex) Freeze() { v.active = false }

// Value returns the value associated with this vertex.
func (v *Vertex) Value() any { return v.value }

// SetValue sets the value associated with this vertex.
func (v *Vertex) SetValue(val any) { v.value = val }

// Edge represents a directed edge in the Graph.
type Edge struct {
	value any
	dstID string
}

// DstID returns the vertex ID that corresponds to this edge's target endpoint.
func (e *Edge) DstID() string { return e.dstID }

// Value returns the value associated with this edge.
func (e *Edge) Value() any { return e.value }

// SetValue sets the value associated with this edge.
func (e *Edge) SetValue(val any) { e.value = val }

