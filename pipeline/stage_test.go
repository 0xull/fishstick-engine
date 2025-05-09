package pipeline_test

// import (
// 	"context"

// 	"github.com/IkehAkinyemi/fishstick-engine/pipeline"
// 	gc "gopkg.in/check.v1"
// )

// var _ = gc.Suite(new(StageTestSuite))

// type StageTestSuite struct {}

// func (s StageTestSuite) TestFIFO(c *gc.C) {
// 	stages := make([]pipeline.StageRunner, 10)
// 	for i := 0; i < len(stages); i++ {
// 		stages[i] = pipeline.FIFO(makePassthroughProcessor())
// 	}
	
// 	src := &sourceStub{data: stringPayloads(3)}
// 	sink := new(sinkStub)
	
// 	p := pipeline.New(stages...)
// 	err := p.Process(context.TODO(), src, sink)
// 	c.Assert(err, gc.IsNil)
// 	c.Assert(sink.data, gc.DeepEquals, src.data)
// 	assertAllProcess(c, src.data)
// }

// func (s StageTestSuite) TestFixedWorkerPool(c *gc.C) {
// 	numWorkers := 10
// 	syncCh := make(chan struct{})
	
// 	rendezvousCh := make(chan struct{})
	
// 	proc := pipeline.ProcessorFunc(func(_ context.Context, _ pipeline.Pipeline) (pipeline.Payload, error)) {
// 		// Signal that we have reached the sync point and wait for the
// 		// green light to proceed by the test code.
// 		syncCh <-struct{}{}
// 		<-rendezvousCh
// 		return nil, nil
// 	}
	
// 	src := &sourceStub{data: stringPayloads(numWorkers)}

// 	p := pipeline.New(pipeline.FixedWorkerPool(proc, numWorkers))
// 	doneCh := make(chan struct{})
// 	go func() {
// 		err := p.Process(context.TODO(), src, nil)
// 		c.Assert(err, gc.IsNil)
// 		close(doneCh)
// 	}()
	
// 	// Wait for all workers to reach sync point. This means that each input
// 	// from the source is currently handled by worker in parallel.
// }
