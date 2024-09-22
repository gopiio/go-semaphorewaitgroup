package semaphorewaitgroup

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreWaitGroup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxConcurrent := int64(2)
	swg := NewSemaphoreWaitGroup(ctx, maxConcurrent)

	// track the concurrency
	concurrent := int64(0)
	maxObservedConcurrent := int64(0)
	mu := sync.Mutex{}

	// A function to simulate the concurrent worker
	doWork := func() {
		defer swg.Done()

		mu.Lock()
		concurrent++
		if concurrent > maxObservedConcurrent {
			maxObservedConcurrent = concurrent
		}
		mu.Unlock()

		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		concurrent--
		mu.Unlock()
	}

	for i := 0; i < 5; i++ {
		swg.Add(1)
		go doWork()
	}
	swg.Wait()

	// maximum observed concurrent tasks does not exceed the limit
	if maxObservedConcurrent > maxConcurrent {
		t.Errorf("Expected max concurrent tasks to be <= %d, but got %d", maxConcurrent, maxObservedConcurrent)
	}

	// Check that all tasks are done
	if concurrent != 0 {
		t.Errorf("Expected all tasks to be done, but %d are still running", concurrent)
	}

	// should do nothing
	swg.Add(-1)

	// Ensure it doesn't change the behavior of the WaitGroup
	if maxObservedConcurrent != 2 {
		t.Errorf("Expected max observed concurrent to remain 2, but got %d", maxObservedConcurrent)
	}
}

func TestSemaphoreWaitGroupNegativeDelta(t *testing.T) {
	ctx := context.Background()
	swg := NewSemaphoreWaitGroup(ctx, 2)
	swg.Add(-1)
	swg.Wait()
}

func TestSemaphoreWaitGroupContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	swg := NewSemaphoreWaitGroup(ctx, 2)
	swg.Add(1)
	go func() {
		defer swg.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	cancel()
	swg.Wait()
}
