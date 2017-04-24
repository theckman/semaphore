# semaphore
[![License](https://img.shields.io/github/license/theckman/semaphore.svg)](https://github.com/theckman/semaphore/blob/master/LICENSE.txt)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/theckman/semaphore)
[![Latest Git Tag](https://img.shields.io/github/tag/theckman/semaphore.svg)](https://github.com/theckman/semaphore/releases)
[![Travis master Build Status](https://img.shields.io/travis/theckman/semaphore/master.svg?label=TravisCI)](https://travis-ci.org/theckman/semaphore/branches)
[![CircleCI master Build Status](https://img.shields.io/circleci/project/github/theckman/semaphore/master.svg?label=CircleCI)](https://circleci.com/gh/theckman/semaphore)
[![Shippable master Build Status](https://img.shields.io/shippable/58fc497dbaa5e307002c66f5/master.svg?label=Shippable)](https://app.shippable.com/github/theckman/semaphore/status/dashboard)

Package `semaphore` is a Go library that implements a semaphore for controlling
the concurrency of workloads. Using a buffered channel, this package provides a
safe way for multiple goroutines to request and return semaphore permits.

For those not familiar with the concept of a semaphore, it allows you to control
the access to a particular resource (or particular number of resources). A
common use case is to limit the number of concurrent threads permitted to
process work, to avoid overcomitting the resources of a system. More details
(historical and applied) about semaphores can be found
[here](https://en.wikipedia.org/wiki/Semaphore_(programming)).

## License

This package is licensed under the permissive [MIT License](https://github.com/theckman/semaphore/blob/master/LICENSE.txt).

## Contributions

Contributions back to this project are more than welcomed. To request a change
please issue a pull request from a fork of this repo, and be sure to include
details about the change in your commit message.

## Installation

Until `go dep` is in wider utilization:

```Shell
go get -u github.com/theckman/semaphore
```

## Usage

The package uses three method for working with semaphores: `Acquire()`,
`Release()`, and `Close()`. Detailed usage information can be discovered via
[GoDoc](https://godoc.org/github.com/theckman/semaphore).

Here is a quick overview of using this package:

```Go
import "github.com/theckman/semaphore"

/* ... */

sema, err := semaphore.New(2)
if err != nil {
	// handle error  
}

for i := 0; i < 10; i++ {
	// blocks until semaphore permit is given
	// or until semaphore has Close() called
	err := sema.Acquire()
	if err != nil {
		// lock acquistion failed (i.e., semaphore is no longer usable)
		panic("semaphore in unexpected state")
	}
    
	// lock acquired, spin-off work
	go func() {
		time.Sleep(time.Second * 3)
		sema.Release()
	}()
}

/* wait for work to finish... */

sema.Close()
```
