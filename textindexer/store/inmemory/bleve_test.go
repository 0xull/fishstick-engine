package inmemory

import (
	"testing"

	"github.com/IkehAkinyemi/fishstick-engine/textindexer/index/indextest"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(InMemoryBleveIndexerTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type InMemoryBleveIndexerTestSuite struct {
	indextest.SuiteBase
	
	idx *InMemoryBleveIndexer
}

func (s *InMemoryBleveIndexerTestSuite) SetupTest(c *gc.C) {
	idx, err := NewInMemoryBleveIndexer()
	c.Assert(err, gc.IsNil)
	s.SetIndexer(idx)
	s.idx = idx
}

func (s *InMemoryBleveIndexerTestSuite) TearDownTest(c *gc.C) {
	c.Assert(s.idx.Close(), gc.IsNil)
}