package pipeline

import (
	"context"
	"fmt"
	"sync"
)

type ProcessStatus int

const (
	// Send the output to the channel
	ProcessOutput ProcessStatus = iota
	// Skip sending the output to the channel, but keep processing
	ProcessSkip
	// Stop processing altogether
	ProcessAbort
)

// Process the input into an output.
//
// The processor should check context.Context in order to know when to early stop its work,
// for instance use http.NewRequestWithContext() or check context.Context.Done() explicitly.
//
// The context also contains a logger which can be retrieved from the context.
type Processor[I any, O any] func(ctx context.Context, input I) (output O, status ProcessStatus)

// Core loop over a channel, used by both Process and ParallelProcess
func processChannel[I any, O any](
	ctx context.Context,
	inputChan <-chan I,
	outputChan chan<- O,
	processor Processor[I, O],
) (finishedInput bool) {
	logger := FromContext(ctx)

	for input := range inputChan {
		select {
		case <-ctx.Done():
			logger.Debugw("context.Done()", "err", ctx.Err())
			return
		default:
		}

		output, status := processor(ctx, input)
		switch status {
		case ProcessSkip:
			logger.Debugw("skip", "input", input, "output", output)
			continue

		case ProcessAbort:
			logger.Debugw("abort", "input", input, "output", output)
			return

		case ProcessOutput:
			select {
			case <-ctx.Done():
				logger.Debugw("context.Done()", "err", ctx.Err())
				return

			case outputChan <- output:
				logger.Debugw("processed", "input", input, "output", output)
			}
		}
	}
	logger.Debug("input channel ended")
	finishedInput = true
	return
}

func finalizeProcess[O any](
	ctx context.Context,
	outputChan chan<- O,
	finalize Finalize[O],
	finishedInput bool,
) {
	if finalize == nil {
		return
	}
	logger := FromContext(ctx)
	output, status := finalize(ctx, finishedInput)
	logger.Debugw("finalized", "status", status, "output", output)
	if status != ProcessOutput {
		return
	}

	select {
	case <-ctx.Done():
		logger.Debugw("context.Done()", "err", ctx.Err())
	case outputChan <- output:
	}
}

// Finalize processing the input.
//
// If finishedInput is true, the input channel was exhausted/closed and all input items were processed.
// If it's false, it was early aborted by context.Context.Done() or ProcessAbort status.
//
// The context also contains a logger which can be retrieved from the context.
type Finalize[O any] func(ctx context.Context, finishedInput bool) (output O, status ProcessStatus)

// General framework to process items from one channel to another.
//
// Processing is done to an un-buffered output channel, one by one, so it will
// block until the next item can be consumed.
//
// Processing may be early stopped by context.Context.Done(), see
// context.WithCancel(), context.WithTimeout() and context.WithDeadline()
//
// The function finalize may be provided and is called after all the processing is done and before
// the channel is closed.
//
// The context may also contain a logger which can be retrieved with ContextLogger(),
// which defaults to the module logger. A new logger will be derived from it including
// both inputChan and outputChan values and this new logger will be passed to the processor.
func Process[I any, O any](
	ctx context.Context,
	inputChan <-chan I,
	processor Processor[I, O],
	finalize Finalize[O],
) (outputChan <-chan O) {
	ch := make(chan O)
	outputChan = ch

	logger := FromContext(ctx).Named("Process").With(
		"inputChan", fmt.Sprintf("%#v", inputChan),
		"outputChan", fmt.Sprintf("%#v", outputChan),
	)
	ctx = NewContext(ctx, logger)

	generator := func() {
		defer func() {
			logger.Info("closing output channel")
			close(ch)
		}()
		finishedInput := processChannel[I, O](ctx, inputChan, ch, processor)
		finalizeProcess(ctx, ch, finalize, finishedInput)
	}

	logger.Info("start")
	go generator()
	return
}

// Reads batchSize entries from the input channel and send them as an slice <= batchSize into the output channel
//
// Batching is done on an un-buffered channel, one by one, so it will
// block until the next batch can be consumed.
//
// Batching may be early stopped by context.Context.Done(), see
// context.WithCancel(), context.WithTimeout() and context.WithDeadline()
func Batch[T any](
	ctx context.Context,
	inputChan <-chan T,
	batchSize int,
) <-chan []T {
	logger := FromContext(ctx).Named("Batch").With(
		"batchSize", batchSize,
	)
	ctx = NewContext(ctx, logger)

	var batch []T

	return Process[T, []T](ctx, inputChan, func(ctx context.Context, input T) (output []T, status ProcessStatus) {
		if batch == nil {
			batch = make([]T, 0, batchSize)
		}
		batch = append(batch, input)
		if len(batch) < cap(batch) {
			return nil, ProcessSkip
		}

		output = batch
		batch = nil
		return output, ProcessOutput
	}, func(ctx context.Context, finishedInput bool) (output []T, status ProcessStatus) {
		if finishedInput && batch != nil {
			return batch, ProcessOutput
		}
		return nil, ProcessSkip
	})
}

// Helper to produce output that is paired with input and error
type ProcessorResult[I any, O any] struct {
	Input  I
	Output O
	Err    error
}

// Process items in parallel (fan-out) with a maximum number of parallel workers.
//
// Note that each item may be processed in a different goroutine.
//
// Note that processor() operations using io.Reader/Read() are **NOT** thread safe as
// they will advance the offset/pointer. However only one item is expected to be processed
// at time, so there should be no worries.
//
// Processing may be early stopped by context.Context.Done(), see
// context.WithCancel(), context.WithTimeout() and context.WithDeadline().
// If any processor() returns an error the processing is **NOT** early stopped, one
// must cancel the context explicitly.
//
// The processor() should check context.Context in order to know when to early stop its work,
// for instance use http.NewRequestWithContext() or check context.Context.Done() explicitly.
//
// Finalize is called with finishedInput == true if all parallel workers finished their input
// and false if at least one of them stopped earlier.
//
// Code is based on https://go.dev/blog/pipelines#bounded-parallelism
func ParallelProcess[I any, O any](
	ctx context.Context,
	maxParallelProcessors int,
	inputChan <-chan I,
	processor Processor[I, O],
	finalize Finalize[O],
) (outputChan <-chan O) {
	ch := make(chan O)

	logger := FromContext(ctx).Named("ParallelProcess").With(
		"maxParallelProcessors", maxParallelProcessors,
		"inputChan", fmt.Sprintf("%#v", inputChan),
		"outputChan", fmt.Sprintf("%#v", ch),
	)
	ctx = NewContext(ctx, logger)

	finishedInput := true

	wg := sync.WaitGroup{}
	wg.Add(maxParallelProcessors)
	innerLogger := FromContext(ctx)
	for i := 0; i < maxParallelProcessors; i++ {
		workerCtx := NewContext(ctx, innerLogger.With("worker", i))
		go func() {
			fi := processChannel[I, O](workerCtx, inputChan, ch, processor)
			if !fi {
				finishedInput = false
			}
			wg.Done()
		}()
	}

	go func() {
		defer func() {
			logger.Info("closing output channel")
			close(ch)
		}()
		wg.Wait()
		finalizeProcess(ctx, ch, finalize, finishedInput)
	}()

	return ch
}
