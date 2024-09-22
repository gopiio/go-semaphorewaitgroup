// Package semaphorewaitgroup provide semaphore to standard WaitGroup.
// This control the number of goroutines can run on a given time
package semaphorewaitgroup

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type SemaphoreWaitGroup struct {
	ctx context.Context
	wg  sync.WaitGroup
	sem *semaphore.Weighted
}

// NewSemaphoreWaitGroup creates new semaphore waitgroup
func NewSemaphoreWaitGroup(ctx context.Context, max int64) *SemaphoreWaitGroup {
	return &SemaphoreWaitGroup{
		ctx: ctx,
		sem: semaphore.NewWeighted(max),
	}
}

func (swg *SemaphoreWaitGroup) Add(delta int) {
	if delta < 0 {
		return
	}

	swg.sem.Acquire(swg.ctx, int64(delta))
	swg.wg.Add(delta)
}

func (swg *SemaphoreWaitGroup) Done() {
	swg.sem.Release(1)
	swg.wg.Done()
}

func (swg *SemaphoreWaitGroup) Wait() {
	swg.wg.Wait()
}
