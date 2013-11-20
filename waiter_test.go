package waiter

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	w := New(2)

	go func() {
		w.Done <- true
	}()

	go func() {
		w.Wrap(nil)
	}()

	if err := w.Wait(); err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	w := New(1)

	go func() {
		w.Errors <- errors.New("an error")
	}()

	if err := w.Wait(); !reflect.DeepEqual(err, errors.New("an error")) {
		t.Error("Errors not right", err)
	}
}

type T bool

func (t T) Close() error {
	if t {
		return nil
	} else {
		return errors.New("err")
	}
}

func TestCloser(t *testing.T) {
	w := New(2)

	w.Close(T(true))
	w.Close(T(false))

	if err := w.Wait(); !reflect.DeepEqual(err, errors.New("err")) {
		t.Error("Error isn't right", err)
	}
}

func TestBypass(t *testing.T) {

	w := New(2)

	w.Close(T(true))
	w.Bypass()

	if err := w.Wait(); err != nil {
		t.Error(err)
	}
}

func TestMultipleErrors(t *testing.T) {
	w := New(2)

	go func() {
		w.Errors <- errors.New("something")
		w.Done <- true
	}()

	go func() {
		w.Errors <- errors.New("else")
		w.Done <- true
	}()

	if err := w.Wait(); err == nil {
		t.Error("Err is nil?")
	}
}

func TestNilCase(t *testing.T) {
	w := New(0)
	if err := w.Wait(); err != nil {
		t.Error(err)
	}
}

func TestCloseNil(t *testing.T) {
	w := New(1)
	var c1 *T
	w.Close(c1)
	if err := w.Wait(); err != nil {
		t.Error(err)
	}
}

func TestReturnTime(t *testing.T) {
	w := New(2)
	begin := time.Now()
	go func() {
		time.Sleep(10 * time.Millisecond)
		w.Done <- true
	}()

	go func() {
		w.Done <- true
	}()

	if err := w.Wait(); err != nil {
		t.Error(err)
	} else if time.Now().Sub(begin) < 10*time.Millisecond {
		t.Error("Time took too little", time.Now().Sub(begin))
	}
}

func TestErrorReturnTime(t *testing.T) {
	w := New(2)
	begin := time.Now()
	go func() {
		time.Sleep(10 * time.Millisecond)
		w.Done <- true
	}()

	go func() {
		w.Errors <- errors.New("hello")
		w.Done <- true
	}()

	if err := w.Wait(); err == nil {
		t.Error("Error is nil?")
	} else if time.Now().Sub(begin) > 10*time.Millisecond {
		t.Error("Time took too long", time.Now().Sub(begin))
	}
}

func TestTimeout(t *testing.T) {
	w := New(1)
	begin := time.Now()
	if err := w.WaitTimeout(50 * time.Millisecond); err == nil {
		t.Error("Error is not timeout?")
	}
	if time.Now().Sub(begin) < 50*time.Millisecond {
		t.Error("Timeout took too short!")
	}
}

func TestTimeout2(t *testing.T) {
	w := New(1)
	begin := time.Now()
	w.Close(nil)
	if err := w.WaitTimeout(50 * time.Millisecond); err != nil {
		t.Error("Error shouldn't be returned", err)
	}
	if time.Now().Sub(begin) > 50*time.Millisecond {
		t.Error("Shouldn't have waited until timeout!")
	}
}
