package inmemory

import (
	"testing"

	"github.com/IkehAkinyemi/fishstick-engine/linkgraph/graph/graphtest"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(InMemoryGraphTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type InMemoryGraphTestSuite struct {
	graphtest.SuiteBase
}

func (s *InMemoryGraphTestSuite) SetupTest(c *gc.C) {
	s.SetGraph(NewInMemoryGraph())
}
