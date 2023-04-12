package concurrent

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
)

// BatchFunc is called for each batch.
// Any error will cancel the batching operation but returning Abort
// indicates it was deliberate, and not an error case.
// [start:end)
type BatchFunc func(start, end int) error

type ParallelBatchFunc func(ctx context.Context, start, end int) error

// Abort indicates a batch operation should abort early.
var (
	Abort      = errors.New("done")
	defaultErr = errors.New("default")
	nilErr     = errors.New("nil")
)

// All calls eachFn for all items
// Returns any error from eachFn except for Abort it returns nil.
func All(count, batchSize int, eachFn BatchFunc) error {
	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}
		err := eachFn(i, end)
		if err == Abort {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func Parallel(ctx context.Context, count, batchSize, maxWorker int, eachFn ParallelBatchFunc) error {
	if count == 0 {
		return nil
	}
	maxWorker, totalInvokeCount, err := splitParallel(count, batchSize, maxWorker)
	if err != nil {
		return err
	}
	return parallelInternal(ctx, count, maxWorker, totalInvokeCount, batchSize, eachFn)
}

func splitParallel(count, batchSize, maxWorker int) (int, int, error) {
	if count == 0 {
		return 0, 0, nil
	}
	if batchSize <= 0 {
		return 0, 0, errors.New("parallel batch size cannot be 0")
	}
	if maxWorker <= 0 {
		return 0, 0, errors.New("max worker cannot be 0")
	}

	totalInvokeCount := count / batchSize // how many times eachFn is called
	if count%batchSize != 0 {             // have leftovers, increase totalInvokeCount by 1
		totalInvokeCount += 1
	}
	if maxWorker > totalInvokeCount { // some workers will not get a job, let's make maxWorker = totalInvokeCount
		maxWorker = totalInvokeCount
	}

	return maxWorker, totalInvokeCount, nil
}

func parallelInternal(ctx context.Context, count, maxWorker, totalInvokeCount, batchSize int, eachFn ParallelBatchFunc) error {
	dataChan := make(chan int, totalInvokeCount)
	errChan := make(chan error, totalInvokeCount)
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for worker := 0; worker < maxWorker; worker += 1 {
		go func() {
			err := defaultErr
			defer func() {
				r := recover()
				if r != nil {
					errChan <- errors.New(fmt.Sprintf("parallel batch panic recovered: %v, %s", r, string(debug.Stack())))
					return
				}
				//To handle the runtime.Goexit() of gomock
				//This means the eachFn exit and doesn't throw any panic
				if err == defaultErr {
					errChan <- errors.New(fmt.Sprintf("goroutine exited without panic, %s", string(debug.Stack())))
				}
			}()
			for {
				select {
				case <-cancelCtx.Done():
					err = nilErr
					return
				case start, hasJob := <-dataChan:
					if !hasJob {
						err = nilErr
						return
					}
					end := start + batchSize
					if end > count {
						end = count
					}
					err = eachFn(cancelCtx, start, end)
					errChan <- err
				}
			}
		}()
	}

	for i := 0; i < count; i += batchSize {
		dataChan <- i
	}
	close(dataChan)

	finishCount := 0
	for err := range errChan {
		if err != nil {
			if err == Abort {
				return nil
			}
			return err
		}
		finishCount++
		if finishCount >= totalInvokeCount {
			return nil
		}
	}
	return nil
}
