waiter [![Build Status][1]][2]
======

A simple way to join return-less concurrently executing goroutines.

## Documentation

[Documentation](http://godoc.org/yanatan16/gowaiter)

## Install

```
go get github.com/yanatan16/gowaiter
```

## Use

```go
import (
  "github.com/yanatan16/gowaiter"
)

func DoThreeThingsInParallel() error {
  w := waiter.New(3)

  go func () {
    if err := Something(); err != nil {
      w.Errors <- err
    }
    if err := Something2(): err != nil {
      w.Errors <- err
    }
    w.Done <- true
  }()

  go func () {
    // For simple functions, we can use a simpler syntax
    w.Wrap(DoSomethingElse())
  }()

  // For closing a type, theres even simpler syntax
  closer := CloserType(true)
  w.Close(closer)

  // Returns first error or nil after 3 Done's.
  return w.Wait()
}

type CloserType bool
func (t CloserType) Close() error {...}

func DoSomethingElse() error {...}
```

## License

MIT License found in LICENSE file.


[1]: https://travis-ci.org/yanatan16/gowaiter.png?branch=master
[2]: http://travis-ci.org/yanatan16/gowaiter