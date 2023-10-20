# pkgx

> golang public library

[![Go Report Card](https://goreportcard.com/badge/github.com/colinrs/pkgx)](https://goreportcard.com/report/github.com/colinrs/pkgx)
[![Build Status](https://travis-ci.com/colinrs/pkgx.svg?branch=master)](https://travis-ci.com/colinrs/pkgx)
[![GoDoc](https://godoc.org/github.com/colinrs/pkgx?status.svg)](https://godoc.org/github.com/colinrs/pkgx)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/colinrs/pkgx/blob/master/LICENSE)

## Overview

A golang public library that reduces code writing in daily development and improves development efficiency

| package | desc |note|
|--|--|--|
| concurrent | Provides the function of concurrency limitation ||
| contextx | extended context ||
| di | Function of dependency injection ||
| fx | golang concurrent execution, similar to python's gevent.spawn method use ||
| http | http client ||
| logger | logger ||
| shutdown | Graceful shutdown of services ||
| structx | Some structures, such as set, safe_map, stack, time_wheel ||
| utils | Common tools, such as json, copy, ping, cmd execution and other functions ||

## Features

- add kq： 
  - Provides an MQ consumer/producer package. Can provide concurrent consumption
  - sequential submission, current limiting, circuit breaker and other functions

## Installation

```shell
go get -u github.com/colinrs/pkgx
```

# Usage Example

```go

package main

import (
  "github.com/colinrs/pkgx/logger"
)

func main() {
  logger.Info("test")
}
```
