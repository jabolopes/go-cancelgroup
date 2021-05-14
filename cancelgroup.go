package cancelgroup

import (
	"context"
	"sync"
)

// Group is a sync.WaitGroup with support for context cancellation.
//
// Goroutines can be scheduled on a Group. The scheduled goroutines are a passed
// a context that is derived from the context passed in the Group's constructor,
// and which is cancelled when the group is cancelled.
//
// Goroutines can be freely scheduled and waited upon. Once the Group is
// scheduled, further goroutines won't be scheduled.
//
// This is safe for concurrent use.
type Group struct {
	ctx        context.Context
	cancel     func()
	cancelOnce sync.Once
	group      sync.WaitGroup
}

// Go schedules a new goroutine in this Group if the Group is not cancelled,
// otherwise this is no-op.
func (g *Group) Go(f func(context.Context)) *Group {
	if g.ctx.Err() != nil {
		return g
	}

	g.group.Add(1)
	go func() {
		defer g.group.Done()
		f(g.ctx)
	}()

	return g
}

// Cancel marks this Group as cancelled. Goroutines that are already scheduled
// will continue running until they complete. Calling cancel multiple times is
// safe but has no further effect once the Group is already cancelled.
func (g *Group) Cancel() *Group {
	g.cancelOnce.Do(g.cancel)
	return g
}

// Wait waits for all scheduled goroutines to complete.
func (g *Group) Wait() *Group {
	g.group.Wait()
	return g
}

// New creates a new Group. The goroutine context is derived from the given
// context.
func New(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)

	group := &Group{
		ctx,
		cancel,
		sync.Once{},
		sync.WaitGroup{},
	}

	return group
}
