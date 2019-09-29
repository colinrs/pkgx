package concurrent

import (
	"runtime"
	"sync"
	"testing"
)

func TestLimit(t *testing.T) {
	limit := NewLimit(2)
	if !(limit.TryAcquire() && limit.TryAcquire() && !limit.TryAcquire()) {
		t.Error("error, TryAcquire")
	}

	limit.Release()
	limit.Release()
}

func BenchmarkLimit(b *testing.B) {
	b.StopTimer()
	limit := NewLimit(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if limit.TryAcquire() {
			limit.Release()
		}
	}
}

func BenchmarkLimitConcurrent(b *testing.B) {
	b.StopTimer()
	limit := NewLimit(1)
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	each := b.N / workers
	wg.Add(workers)
	b.StartTimer()
	for i := 0; i < workers; i++ {
		go func() {
			for i := 0; i < each; i++ {
				if limit.TryAcquire() {
					limit.Release()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
