// Copyright (c) 2017 Tim Heckman
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE.txt file.

// Package semaphore is a Go library that implements a semaphore for controlling
// the concurrency of workloads. Using a buffered channel, this package provides a
// safe way for multiple goroutines to request and return semaphore permits.
//
// For those not familiar with the concept of a semaphore, it allows you to control
// the access to a particular resource (or particular number of resources). A
// common use case is to limit the number of concurrent threads permitted to
// process work, to avoid overcomitting the resources of a system. More details
// (historical and applied) about semaphores can be found
// here: https://en.wikipedia.org/wiki/Semaphore_(programming).
//
// The package uses three method for working with semaphores: Acquire(),
// Release(), and Close(). Here is a quick overview of using this package:
//
// 		import "github.com/theckman/semaphore"
//
// 		/* ... */
//
// 		sema, err := semaphore.New(2)
// 		if err != nil {
// 			// handle error
// 		}
//
// 		for i := 0; i < 10; i++ {
// 			// blocks until semaphore permit is given
// 			// or until semaphore has Close() called
// 			err := sema.Acquire()
// 			if err != nil {
// 				// lock acquistion failed (i.e., semaphore is no longer usable)
// 				panic("semaphore in unexpected state")
// 			}
//
// 			// lock acquired, spin-off work
// 			go func() {
// 				time.Sleep(time.Second * 3)
// 				sema.Release()
// 			}()
// 		}
//
// 		/* wait for work to finish... */
//
// 		sema.Close()
package semaphore
