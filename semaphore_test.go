// Copyright (c) 2017 Tim Heckman
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE.txt file.

package semaphore

import (
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestNew_InvalidSizeArgument(t *testing.T) {
	var sema Semaphore
	var err error
	sema, err = New(0)

	if sema != nil {
		sema.Close()
		t.Fatalf("New(0) semaphore %#v, want <nil>", sema)
	}

	if err == nil {
		t.Fatalf("New(0) error = <nil>, want %q", "size argument must be greater than 0")
	}
}

func TestNew(t *testing.T) {
	var sema Semaphore
	var err error

	sema, err = New(42)

	if err != nil {
		t.Fatalf("New(42) error = %q, want <nil>", err)
	}

	if sema == nil {
		t.Fatalf("New(42) semaphore = <nil>, want != <nil>")
	}

	defer sema.Close()

	castSem, ok := sema.(*semaphore)

	if !ok {
		t.Fatalf("type assertion failed: sem = %s, want *semaphore", reflect.ValueOf(sema).String())
	}

	if semCap := cap(castSem.c); semCap != 42 {
		t.Fatalf("cap(castSem.c) = %d, want 42", semCap)
	}
}

func Test_semaphoreAcquire(t *testing.T) {
	sema := &semaphore{c: make(chan struct{}, 1)}

	defer sema.Close()

	if err := sema.Acquire(); err != nil {
		t.Fatalf("sema.Acquire() = %s, want <nil>", err)
	}

	select {
	case _, ok := <-sema.c:
		if !ok {
			t.Fatal("sema.c should not be closed")
		}
	default:
		t.Fatal("channel is expected to have at least one message")
	}
}

func Test_semaphoreAcquire_Closed(t *testing.T) {
	sema := &semaphore{c: make(chan struct{}, 1)}

	sema.Close()

	if err := sema.Acquire(); err != ErrUnusable {
		t.Fatalf("sema.Acquire() = %s, want %s", err, ErrUnusable)
	}
}

func releaseIt(sema Semaphore, err chan<- error) {
	err <- sema.Release()
}

func Test_semaphoreRelease(t *testing.T) {
	sema := &semaphore{c: make(chan struct{}, 1)}

	defer sema.Close()

	sema.c <- struct{}{}

	errCh := make(chan error)
	timer := time.NewTimer(time.Second * 2)

	go releaseIt(sema, errCh)
	runtime.Gosched()

	select {
	case err := <-errCh:
		timer.Stop()
		close(errCh)
		if err != nil {
			t.Fatalf("sema.Release() = %s, want <nil>", err)
		}
	case _ = <-timer.C:
		t.Fatal("semaphore Release did not complete after 2 seconds")
	}
}

func Test_semaphoreRelease_Closed(t *testing.T) {
	sema := &semaphore{c: make(chan struct{}, 1)}

	sema.c <- struct{}{}

	sema.Close()

	errCh := make(chan error)
	timer := time.NewTimer(time.Second * 2)

	go releaseIt(sema, errCh)
	runtime.Gosched()

	select {
	case err := <-errCh:
		timer.Stop()
		close(errCh)
		if err != nil {
			t.Fatalf("sema.Release() = %s, want <nil>", err)
		}
	case _ = <-timer.C:
		t.Fatal("semaphore Release did not complete after 2 seconds")
	}

}

func encapsuateClose(ch chan struct{}) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	close(ch)

	return
}

func Test_semaphoreClose(t *testing.T) {
	var err error
	sema := &semaphore{c: make(chan struct{})}

	if err = sema.Close(); err != nil {
		t.Fatalf("sema.Close() = %#v, want <nil>", err)
	}

	if panicked := encapsuateClose(sema.c); panicked != true {
		t.Fatal("sema.Close() does not appear to have closed the channel")
	}
}

func Test_semaphoreClose_AlreadyClosed(t *testing.T) {
	var err error
	sema := &semaphore{c: make(chan struct{})}

	if err = sema.Close(); err != nil {
		t.Fatalf("sema.Close() = %#v, want <nil>", err)
	}

	if err = sema.Close(); err != ErrAlreadyClosed {
		t.Fatalf("sema.Close() = %#v, want %#v", err, ErrAlreadyClosed)
	}
}
