// Package waiter provides a nice way to perform return-less parallel operations
//
// Example:
//
//    func DoThreeThingsParallel() error {
//      w := waiter.New(3)
//
//      go func () {
//        // do something
//        if err != nil {
//          w.Errors <- err
//        }
//        // done
//        w.Done <- true
//      }()
//
//      go func () {
//        err := doSomethingElse()
//        w.Wrap(err) // Takes care of errors and done channels
//      }()
//
//      w.Close(CloserInterface)
//
//      return w.wait()
//    }
package waiter

import (
	"io"
	"reflect"
)

type Waiter struct {
	Errors chan error
	Done   chan bool
	Count  int
}

func New(cnt int) *Waiter {
	return &Waiter{make(chan error), make(chan bool), cnt}
}

func (w Waiter) Wait() error {
	for i := 0; i < w.Count; i++ {
		select {
		case err := <-w.Errors:
			w.waitAndClose(i)
			return err
		case <-w.Done:
		}
	}
	close(w.Errors)
	close(w.Done)
	return nil
}

func (w Waiter) Wrap(err error) {
	if err != nil {
		w.Errors <- err
	} else {
		w.Done <- true
	}
}

func (w Waiter) Close(closer io.Closer) {
	go func() {
		if isNotNil(closer) {
			w.Wrap(closer.Close())
		} else {
			w.Done <- true
		}
	}()
}

// Bypass a single wait. That is, we decrement count.
func (w *Waiter) Bypass() {
	w.Count -= 1
}

// After an error we return from wait, but must ensure Done and Errors are properly closed
func (w Waiter) waitAndClose(i int) {
	go func() {
		for j := i; j < w.Count; j++ {
			<-w.Done
		}
		close(w.Done)
		close(w.Errors)
	}()

	go func() {
		for _ = range w.Errors {
		}
	}()
}

func isNotNil(i interface{}) bool {
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).Pointer() != 0
	} else {
		return i != nil
	}
}
