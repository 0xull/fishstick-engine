package bspgraph_test

import (
	"fmt"
	"testing"

	bspgraph "github.com/IkehAkinyemi/fishstick-engine/bsp"
	"github.com/IkehAkinyemi/fishstick-engine/bsp/aggregator"
	"github.com/IkehAkinyemi/fishstick-engine/bsp/message"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(GraphTestSuite))

func Test(t *testing.T) {
	// Run all gocheck test-suites
	gc.TestingT(t)
}

type GraphTestSuite struct{}

func (s *GraphTestSuite) TestMessageExchange(c *gc.C) {
	g, err := bspgraph.NewGraph(bspgraph.GraphConfig{
		ComputeFn: func(g *bspgraph.Graph, v *bspgraph.Vertex, msgIt message.Iterator) error {
			v.Freeze()
			if g.Superstep() == 0 {
				var dst string
				switch v.ID() {
				case "0":
					dst = "1"
				case "1":
					dst = "0"
				}

				return g.SendMessage(dst, &intMsg{value: 42})
			}

			for msgIt.Next() {
				v.SetValue(msgIt.Message().(*intMsg).value)
			}
			return nil
		},
	})
	c.Assert(err, gc.IsNil)
	defer func() { c.Assert(g.Close(), gc.IsNil) }()

	g.AddVertex("0", 0)
	g.AddVertex("1", 0)

	err = execFixedSteps(g, 2)
	c.Assert(err, gc.IsNil)

	for id, v := range g.Vertices() {
		c.Assert(v.Value(), gc.Equals, 42, gc.Commentf("vertex %v", id))
	}
}

func (s *GraphTestSuite) TestMessageBroadcasting(c *gc.C) {
	g, err := bspgraph.NewGraph(bspgraph.GraphConfig{
		ComputeFn: func(g *bspgraph.Graph, v *bspgraph.Vertex, msgIt message.Iterator) error {
			if err := g.BroadcastToNeighbors(v, &intMsg{value: 42}); err != nil {
				return err
			}
			for msgIt.Next() {
				v.SetValue(msgIt.Message().(*intMsg).value)
			}
			return nil
		},
	})
	c.Assert(err, gc.IsNil)
	defer func() { c.Assert(g.Close(), gc.IsNil) }()

	g.AddVertex("0", 42)
	g.AddVertex("1", 0)
	g.AddVertex("2", 0)
	g.AddVertex("3", 0)

	c.Assert(g.AddEdge("0", "1", nil), gc.IsNil)
	c.Assert(g.AddEdge("0", "2", nil), gc.IsNil)
	c.Assert(g.AddEdge("0", "3", nil), gc.IsNil)

	err = execFixedStep(g, 2)
	c.Assert(err, gc.IsNil)

	for id, v := range g.Vertices() {
		c.Assert(v.Value(), gc.Equals, 42, gc.Commentf("vertex %v", id))
	}
}

func (s *GraphTestSuite) TestAggregator(c *gc.C) {
	g, err := bspgraph.NewGraph(bspgraph.GraphConfig{
		ComputeWorkers: 4,
		ComputeFn: func(g *bspgraph.Graph, v *bspgraph.Vertex, msgIt message.Iterator) error {
			g.Aggregator("counter").Aggregate(1)
			return nil
		},
	})
	c.Assert(err, gc.IsNil)
	defer func() { c.Assert(g.Close(), gc.IsNil) }()

	offset := 5
	g.RegisterAggregator("counter", new(aggregator.IntAccumulator))
	g.Aggregator("counter").Aggregate(offset)

	numVerts := 1000
	for i := range numVerts {
		g.AddVertex(fmt.Sprint(i), nil)
	}

	err = execFixedSteps(g, 1)
	c.Assert(err, gc.IsNil)

	aggrMap := g.Aggregators()
	c.Assert(aggrMap["counter"].Get(), gc.Equals, numVerts+offset)
}

func (s *GraphTestSuite) TestMessageRelay(c *gc.C) {
	g1, err := bspgraph.NewGraph(bspgraph.GraphConfig{
		ComputeFn: func(g *bspgraph.Graph, v *bspgraph.Vertex, msgIt message.Iterator) error {
			if g.Superstep() == 0 {
				for _, e := range v.Edges() {
					_ = g.SendMessage(e.DstID(), &intMsg{value: 42})
				}
				return nil
			}

			for msgIt.Next() {
				v.SetValue(msgIt.Message().(*intMsg).value)
			}
			return nil
		},
	})
	c.Assert(err, gc.IsNil)
	defer func() { c.Assert(g1.Close(), gc.IsNil) }()

	g2, err := bspgraph.NewGraph(bspgraph.GraphConfig{
		ComputeFn: func(g *bspgraph.Graph, v *bspgraph.Vertex, msgIt message.Iterator) error {
			for msgIt.Next() {
				m := msgIt.Message().(*intMsg)
				v.SetValue(m.value)
				_ = g.SendMessage("graph1.vertex", m)
			}
			return nil
		},
	})
	c.Assert(err, gc.IsNil)
	defer func() { c.Assert(g2.Close(), gc.IsNil) }()
	
	g2.AddVertex("graph2.vertex", nil)
	g2.RegisterRelayer(localRelayer{to: g1})

	g1.AddVertex("graph1.vertex", nil)
	c.Assert(g1.AddEdge("graph1.vertex", "graph2.vertex", nil), gc.IsNil)
	g1.RegisterRelayer(localRelayer{to: g2})

}
