package cancelgroup_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jabolopes/go-cancelgroup"
)

func ExampleGroup_GroupCancellation() {
	ctx := context.Background()

	g := cancelgroup.New(ctx)
	g.Go(func(ctx context.Context) {
		fmt.Printf("hello")
	})

	time.Sleep(100 * time.Millisecond)

	// Cancel group.
	g.Cancel()

	// Scheduling goroutines on the already cancelled group is a no-op
	g.Go(func(ctx context.Context) {
		fmt.Printf("world")
	})

	// Wait for active goroutines to finish.
	g.Wait()

	// Output: hello
}

func ExampleGroup_ContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())

	g := cancelgroup.New(ctx)
	g.Go(func(ctx context.Context) {
		fmt.Printf("hello")
	})

	time.Sleep(100 * time.Millisecond)

	// Cancel context.
	cancel()

	// Scheduling goroutines on the already cancelled group is a no-op
	g.Go(func(ctx context.Context) {
		fmt.Printf("world")
	})

	// Wait for active goroutines to finish.
	g.Wait()

	// Output: hello
}

func randSleep(i int) {
	time.Sleep(time.Duration(rand.Intn(i+1)) * time.Millisecond)
}

func TestEmptyGroup(t *testing.T) {
	ctx := context.Background()
	_ = cancelgroup.New(ctx)
}

func TestEmptyGroup_Wait(t *testing.T) {
	ctx := context.Background()

	g := cancelgroup.New(ctx)
	g.Wait()
}

func TestGoroutineRuns_Single(t *testing.T) {
	ctx := context.Background()

	done := false

	g := cancelgroup.New(ctx)
	g.Go(func(ctx context.Context) {
		randSleep(10)
		done = true
	})
	g.Wait()

	if !done {
		t.Errorf("done = %v; want %v", done, true)
	}
}

func TestGoroutineRuns_Multiple(t *testing.T) {
	ctx := context.Background()

	const n = 10

	dones := make([]bool, 10)

	g := cancelgroup.New(ctx)

	for i := range dones {
		i := i

		g.Go(func(ctx context.Context) {
			randSleep(i)
			dones[i] = true
		})
	}

	g.Wait()

	if len(dones) != n {
		t.Fatalf("len(dones) = %d; want %v", len(dones), n)
	}

	for i, done := range dones {
		if !done {
			t.Errorf("dones[%d] = %v; want %v", i, done, true)
		}
	}
}

func TestGoroutineCancels_Single(t *testing.T) {
	ctx := context.Background()

	done := false
	doneErr := error(nil)

	g := cancelgroup.New(ctx)
	g.Go(func(ctx context.Context) {
		time.Sleep(100 * time.Millisecond)
		done = true
		doneErr = ctx.Err()
	})

	g.Cancel().Wait()

	if !done || doneErr == nil {
		t.Errorf("done, doneErr = %v, %v; want %v, non-%v", done, doneErr, true, nil)
	}
}

func TestGoroutineCancels_Multiple(t *testing.T) {
	ctx := context.Background()

	const n = 10

	dones := make([]bool, 10)
	doneErrs := make([]error, 10)

	g := cancelgroup.New(ctx)

	for i := range dones {
		i := i

		g.Go(func(ctx context.Context) {
			time.Sleep(100 * time.Millisecond)
			dones[i] = true
			doneErrs[i] = ctx.Err()
		})
	}

	g.Cancel().Wait()

	if len(dones) != n {
		t.Fatalf("len(dones) = %d; want %v", len(dones), n)
	}

	for i := range dones {
		done := dones[i]
		doneErr := doneErrs[i]

		if !done || doneErr == nil {
			t.Errorf("dones[%d] = %v, %v; want %v, non-%v", i, done, doneErr, true, nil)
		}
	}
}

func TestCancelledGroupDoesNotSchedule(t *testing.T) {
	ctx := context.Background()

	done := false

	g := cancelgroup.New(ctx)
	g.Cancel()
	g.Go(func(ctx context.Context) {
		done = true
	})
	g.Wait()

	if done {
		t.Errorf("done = %v; want %v", done, false)
	}
}

func TestCancelledContextDoesNotSchedule(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	done := false

	g := cancelgroup.New(ctx)
	cancel()
	g.Go(func(ctx context.Context) {
		done = true
	})
	g.Wait()

	if done {
		t.Errorf("done = %v; want %v", done, false)
	}
}
