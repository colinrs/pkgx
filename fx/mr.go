package fx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/colinrs/pkgx/contextx"
)

const (
	defaultTimeOut           = 600 * time.Second
	defaultGoConcurrentLimit = 1000
)

var (
	RunFirstErr = errors.New("should run first")
)

type Func func()

type Option func(opts *ConcurrentOptions)

type ConcurrentOptions struct {
	timeout           time.Duration
	goConcurrentLimit int
}

type Concurrent struct {
	waitGroup         RoutineGroup       // wait group
	fgsCount          uint64             // func group count
	fgs               []Func             // func
	result            [][]reflect.Value  // func result
	run               bool               // is run
	concurrentOptions *ConcurrentOptions // options
	err               error
}

func NewConcurrent(options ...Option) *Concurrent {
	concurrent := &Concurrent{
		waitGroup: RoutineGroup{},
		fgsCount:  0,
		concurrentOptions: &ConcurrentOptions{
			timeout:           defaultTimeOut, // default time out
			goConcurrentLimit: defaultGoConcurrentLimit,
		},
	}
	for _, option := range options {
		option(concurrent.concurrentOptions)
	}
	return concurrent
}

// GetParamValues get the running parameters of the function
func (concurrent *Concurrent) GetParamValues(param ...interface{}) []reflect.Value {
	values := make([]reflect.Value, 0, len(param))
	for i := range param {
		values = append(values, reflect.ValueOf(param[i]))
	}
	return values
}

func (concurrent *Concurrent) AddFunc(impl interface{}, paramsValue []reflect.Value) {
	index := concurrent.fgsCount
	ff := func() {
		handlerFuncType := reflect.TypeOf(impl)
		if handlerFuncType.Kind() != reflect.Func {
			panic("not a func")
		}
		f := reflect.ValueOf(impl)
		concurrent.result[index] = f.Call(paramsValue)
	}
	concurrent.fgsCount++
	concurrent.fgs = append(concurrent.fgs, ff)
}

func (concurrent *Concurrent) Run() error {
	concurrent.run = true
	concurrent.result = make([][]reflect.Value, len(concurrent.fgs))
	ctx, cancel := contextx.ShrinkDeadline(context.Background(), concurrent.concurrentOptions.timeout)
	defer cancel()
	done := make(chan struct{})
	fn := func() {
		for _, fg := range concurrent.fgs {
			concurrent.waitGroup.RunGoSafe(fg)
		}
		concurrent.waitGroup.Wait()
		done <- struct{}{}
		close(done)
	}
	GoSafe(fn)

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		concurrent.err = ctx.Err()
		return ctx.Err()
	}
	concurrent.err = nil
	return nil
}

func (concurrent *Concurrent) Result() ([][]interface{}, error) {
	if !concurrent.run {
		fmt.Print("no result, run it first\n")
		return nil, RunFirstErr
	}
	result := make([][]interface{}, len(concurrent.result), len(concurrent.result))
	for i, values := range concurrent.result {
		res := make([]interface{}, len(values), len(values))
		for j, value := range values {
			res[j] = value.Interface()
		}
		result[i] = res
	}
	return result, concurrent.err
}

func WithTimeout(timeout time.Duration) Option {
	return func(ConcurrentOptions *ConcurrentOptions) {
		ConcurrentOptions.timeout = timeout
	}
}

func WithGoConcurrentLimit(limit int) Option {
	return func(ConcurrentOptions *ConcurrentOptions) {
		ConcurrentOptions.goConcurrentLimit = limit
	}
}
