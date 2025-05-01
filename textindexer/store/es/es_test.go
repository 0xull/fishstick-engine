package es

import (
	"os"
	"strings"
	"testing"

	"github.com/IkehAkinyemi/fishstick-engine/textindexer/index/indextest"
	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(ElasticSearchIndexerTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type ElasticSearchIndexerTestSuite struct {
	indextest.SuiteBase
	idx *ElasticSearchIndexer
}

func (s *ElasticSearchIndexerTestSuite) SetupSuite(c *gc.C) {
	nodeList := os.Getenv("ES_NODES")
	if nodeList == "" {
		c.Skip("Missing ES_NODES envvar; skipping elasticsearch-backed index test suite")
	}
	
	idx, err := NewElasticSearchIndexer(strings.Split(nodeList, ","), true)
	c.Assert(err, gc.IsNil)
	s.SetIndexer(idx)
	s.idx = idx
}

func (s *ElasticSearchIndexerTestSuite) SetupTest(c *gc.C) {
	if s.idx.es != nil {
		_, err := s.idx.es.Indices.Delete([]string{indexName})
		c.Assert(err, gc.IsNil)
		err = ensureIndex(s.idx.es)
		c.Assert(err, gc.IsNil)
	}
}