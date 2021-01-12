# go-cancelgroup

This package provides a Group, which is a sync.WaitGroup with support for
context cancellation.

Example:

```go
g := cancelgroup.New(ctx)

// Schedule some goroutines.
for i := 0; i < 10; i++ {
  g.Go(func(ctx context.Context) { ... })
}

// Do something else...

// Cancel group.
g.Cancel()

// Wait for active goroutines to finish.
g.Wait()
```
