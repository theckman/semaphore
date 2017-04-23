package semaphore

import (
	"errors"
	"io"
)

// ErrUnusable is the error returned if the semaphore isn't suitable for use,
// meaning something called the Close() method. If this error is returned, the
// semaphore can no longer issue permits and the New() function must be used to
// allocate a new one.
var ErrUnusable = errors.New("semaphore not suitable for use, please create one with New()")

// ErrAlreadyClosed is the error returned from the Close() function if the
// semaphore has already been shut down. This error should be treated as an
// advisory and NOT a failure.
var ErrAlreadyClosed = errors.New("the semaphore has already had the Close() method called")

// Acquirer is the interface for taking permits from a semaphore.
type Acquirer interface {
	// Acquire is a blocking function to take a permit from the semaphore. If the
	// error returned is nil, the semaphore permit has been acquired and the work
	// can start.
	//
	// If an error is returned, the permit was not given and the work MUST NOT
	// start. This returns ErrUnusable if Close() has been called on the semaphore.
	// In this case you'd need to use New() to obtain a new usable semaphore.
	Acquire() error
}

// Releaser is the interface for releasing permits back to a semaphore.
type Releaser interface {
	// Release is a non-blocking function to release the semaphore. If an error is returned
	// from this function, the release was successful but the semaphore cannot be
	// used to acquire another permit.
	//
	// In other words, if this returns an error it should not be treated as a
	// failure. For example, if Close() is called, followed by Release(), this
	// function would return an ErrUnusable error. In this case you'd need to
	// use New() to obtain a new usable semaphore.
	Release() error
}

// Semaphore is the interface needed to implement a functioning semaphore. One
// function to take permit, one to give back a permit, and a final function to
// tell the semaphore to stop issuing permits.
type Semaphore interface {
	// Acquirer is for taking a semaphore permit. See its comments for more
	// details.
	Acquirer

	// Releaser is for releasing a semaphore permit. See its comments for more
	// details.
	Releaser

	// Close is a non-blocking function that shuts the semaphore down and
	// prevents it from issuing further permits. Consumers would need to create
	// a completely new semaphore, using New(), after calling Close() if they
	// wanted to request another permit.
	//
	// Close terminates the semaphore in a way that ensures any consumer with a
	// permit is able to release that permit without fatal error or deadlock.
	//
	// This function should return an ErrAlreadyClosed error if the semaphore
	// has already been closed, but consumers should treat that as an advisory
	// and not a fatal error.
	io.Closer
}

// New returns a new Semaphore and it takes a size argument to define how many
// concurrent permits the semaphore will issue. The size argument must be
// greater than 0 or an error will be returned.
func New(size int) (Semaphore, error) {
	if size < 1 {
		return nil, errors.New("size argument must be greater than 0")
	}

	return &semaphore{c: make(chan struct{}, size)}, nil
}

type semaphore struct {
	c chan struct{}
}

func (s *semaphore) Acquire() (err error) {
	defer func() {
		// this should catch panics for writing to a closed channel
		if r := recover(); r != nil {
			err = ErrUnusable
		}
	}()

	s.c <- struct{}{}

	return
}

func (s *semaphore) Release() error {
	if _, ok := <-s.c; !ok {
		return ErrUnusable
	}

	return nil
}

func (s *semaphore) Close() (err error) {
	// XXX(theckman): catch panic if semaphore channel already closed,
	// return ErrAlreadyClosed if so
	defer func() {
		if r := recover(); r != nil {
			err = ErrAlreadyClosed
		}
	}()

	close(s.c)

	return
}
