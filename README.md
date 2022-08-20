# go-cancelgroup

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jabolopes/go-cancelgroup)](https://pkg.go.dev/github.com/jabolopes/go-cancelgroup)

This package provides `cancelgroup.Group`, which is a `sync.WaitGroup` with
support for group cancellation and context cancellation. When the context is
cancelled, the `cancelgroup.Group` is also cancelled and attempting to further
schedule goroutines on this group is a no-op.

## Installation

```sh
$ go get github.com/jabolopes/go-cancelgroup
```

You can use `go get -u` to update the package. If you are using Go modules, you
can also just import the package and it will be automatically downloaded on the
first compilation.

## Examples

Typical usage of a `cancelgroup.Group` with group cancellation:

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

Example of a usage of a `cancelgroup.Group` with Context cancellation:

```go
// Create a context with cancellation. Can also be a `context.WithDeadline`, a `context.WithTimeout`, etc.
ctx, cancel := context.WithCancel(ctx)

g := cancelgroup.New(ctx)

// Schedule some goroutines.
for i := 0; i < 10; i++ {
  g.Go(func(ctx context.Context) { ... })
}

// Cancel the context.
cancel()

// Scheduling goroutines on the already cancelled group is a no-op
g.Go(func(ctx context.Context) { ... })

// Wait for active goroutines to finish.
g.Wait()
```
