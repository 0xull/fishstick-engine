package aggregator

import (
	"math"
	"math/rand"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(new(AccumulatorTestSuite))

type aggregator interface {
	Set(any)
	Get() any
	Aggregate(any)
	Delta() any
}

func Test(t *testing.T) {
	// Run all gocheck test-suites
	gc.TestingT(t)
}

type AccumulatorTestSuite struct {}

func (s *AccumulatorTestSuite) TestFloat64Accumulator(c *gc.C) {
	numVal := 100
	vals := make([]any, numVal)
	var exp float64
	for i := range numVal {
		next := rand.Float64()
		vals[i] = next
		exp += next
	}
	
	got := s.testConcurrentAccess(new(Float64Accumulator), vals).(float64)
	absDelta := math.Abs(exp - got)
	c.Assert(
		absDelta < 1e-6, gc.Equals, true, 
		gc.Commentf("expected to get %f; got %f; |delta| %f > 1e-6", exp, got, absDelta),
	)
}

func (s *AccumulatorTestSuite) TestIntAccumulator(c *gc.C) {
	numVal := 100
	vals := make([]any, numVal)
	var exp int 
	for i := range numVal {
		next := rand.Int()
		vals[i] = next
		exp += next
	}
	
	got := s.testConcurrentAccess(new(IntAccumulator), vals).(int)
	c.Assert(got, gc.Equals, exp)
}

func (s *AccumulatorTestSuite) testConcurrentAccess(a aggregator, values []any) any {
	startedCh := make(chan struct{})
	synCh := make(chan struct{})
	doneCh := make(chan struct{})
	for i := range len(values) {
		go func(i int) {
			startedCh <- struct{}{}
			<-synCh
			a.Aggregate(values[i])
			doneCh <- struct{}{}
		}(i)
	}
	
	// Wait for all go-routines to start
	for range len(values) {
		<-startedCh
	}
	
	// Allow each go-routine to update their accumulator
	close(synCh)
	
	// Wait for all go-routine to exit
	for range len(values) {
		<-doneCh
	}
	
	return a.Get()
}
