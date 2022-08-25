package pprof

import (
	"bytes"
	"os"
	"path"
	"runtime/pprof"
	"time"

	gosafe "github.com/colinrs/pkgx/fx"
)

const (
	binaryDump = 0
)

const (
	goroutine = iota
	gcHeap
	thread
	block
	mutex
)

const (
	defaultLoggerFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultLoggerPerm  = 0644
)

// check type to check name
var check2name = map[int]string{
	goroutine: "goroutine",
	gcHeap:    "GCHeap",
	thread:    "thread",
	block:     "block",
	mutex:     "mutex",
}

type typeOption struct {
	pprofType int
	filePath  string
	Duration  time.Duration
}

type OptionFun func(o *Option)

type Option struct {
	goroutineOpts *typeOption
	gCHeapOpts    *typeOption
	threadOpts    *typeOption
	blockOpts     *typeOption
	mutexOpts     *typeOption
}

func ProfileAutoFetch(opts ...OptionFun) {
	option := new(Option)
	for _, opt := range opts {
		opt(option)
	}
	runG(option)
}

func runG(option *Option) {
	if option.goroutineOpts != nil {
		gosafe.GoSafe(func() {
			fetchGoroutinePprof(option.goroutineOpts.filePath, option.goroutineOpts.Duration)
		})
	}

	if option.gCHeapOpts != nil {

		gosafe.GoSafe(func() {
			fetchGCHeapPprof(option.gCHeapOpts.filePath, option.gCHeapOpts.Duration)
		})
	}

	if option.threadOpts != nil {
		gosafe.GoSafe(func() {
			fetchThreadCreatePprof(option.threadOpts.filePath, option.threadOpts.Duration)
		})
	}

	if option.blockOpts != nil {
		gosafe.GoSafe(func() {
			fetchBlockPprof(option.blockOpts.filePath, option.blockOpts.Duration)
		})
	}

	if option.mutexOpts != nil {
		gosafe.GoSafe(func() {
			fetchMutexPprof(option.mutexOpts.filePath, option.mutexOpts.Duration)
		})
	}
}

func WithGoroutineOpts(filePath string, duration time.Duration) OptionFun {
	return func(o *Option) {
		o.goroutineOpts = new(typeOption)
		o.goroutineOpts.pprofType = goroutine
		o.goroutineOpts.filePath = filePath
		o.goroutineOpts.Duration = duration
	}
}

func WithGCHeapOpts(filePath string, duration time.Duration) OptionFun {
	return func(o *Option) {
		o.gCHeapOpts = new(typeOption)
		o.gCHeapOpts.pprofType = gcHeap
		o.gCHeapOpts.filePath = filePath
		o.gCHeapOpts.Duration = duration
	}
}

func WithThreadHeapOpts(filePath string, duration time.Duration) OptionFun {
	return func(o *Option) {
		o.threadOpts = new(typeOption)
		o.threadOpts.pprofType = thread
		o.threadOpts.filePath = filePath
		o.threadOpts.Duration = duration
	}
}

func WithMutexHeapOpts(filePath string, duration time.Duration) OptionFun {
	return func(o *Option) {
		o.mutexOpts = new(typeOption)
		o.mutexOpts.pprofType = mutex
		o.mutexOpts.filePath = filePath
		o.mutexOpts.Duration = duration
	}
}

func WithBlockHeapOpts(filePath string, duration time.Duration) OptionFun {
	return func(o *Option) {
		o.blockOpts = new(typeOption)
		o.blockOpts.pprofType = block
		o.blockOpts.filePath = filePath
		o.blockOpts.Duration = duration
	}
}

func durationRun(duration time.Duration, f func()) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for range ticker.C {
		f()
	}
}

func fetchGoroutinePprof(dumpPath string, duration time.Duration) {
	durationRun(duration, func() {

		_, err := GoroutinePprof(dumpPath)
		if err != nil {
			return
		}

	})
}

func GoroutinePprof(dumpPath string) (string, error) {
	var buf bytes.Buffer
	err := pprof.Lookup("goroutine").WriteTo(&buf, binaryDump)
	if err != nil {
		return "", err
	}
	result, err := writeProfileDataToFile(buf, goroutine, dumpPath)
	if err != nil {
		return "", err
	}
	return result, nil
}

func fetchGCHeapPprof(dumpPath string, duration time.Duration) {

	durationRun(duration, func() {
		_, err := GCHeapPprof(dumpPath)
		if err != nil {
			return
		}
	})
}

func GCHeapPprof(dumpPath string) (string, error) {
	var buf bytes.Buffer
	err := pprof.Lookup("heap").WriteTo(&buf, binaryDump)
	if err != nil {
		return "", err
	}
	result, err := writeProfileDataToFile(buf, gcHeap, dumpPath)
	if err != nil {
		return "", err
	}
	return result, nil
}

func fetchThreadCreatePprof(dumpPath string, duration time.Duration) {
	durationRun(duration, func() {
		_, err := ThreadPprof(dumpPath)
		if err != nil {
			return
		}
	})
}

func ThreadPprof(dumpPath string) (string, error) {
	var buf bytes.Buffer
	err := pprof.Lookup("threadcreate").WriteTo(&buf, binaryDump)
	if err != nil {
		return "", err
	}
	result, err := writeProfileDataToFile(buf, thread, dumpPath)
	if err != nil {
		return "", err
	}
	return result, nil
}

func fetchBlockPprof(dumpPath string, duration time.Duration) {
	durationRun(duration, func() {
		_, err := BlockPprof(dumpPath)
		if err != nil {
			return
		}
	})
}

func BlockPprof(dumpPath string) (string, error) {
	var buf bytes.Buffer
	err := pprof.Lookup("block").WriteTo(&buf, binaryDump)
	if err != nil {
		return "", err
	}
	result, err := writeProfileDataToFile(buf, block, dumpPath)
	if err != nil {
		return "", err
	}
	return result, nil
}

func fetchMutexPprof(dumpPath string, duration time.Duration) {
	durationRun(duration, func() {
		_, err := MutexPprof(dumpPath)
		if err != nil {
			return
		}
	})
}

func MutexPprof(dumpPath string) (string, error) {
	var buf bytes.Buffer
	err := pprof.Lookup("mutex").WriteTo(&buf, binaryDump)
	if err != nil {
		return "", err
	}
	result, err := writeProfileDataToFile(buf, mutex, dumpPath)
	if err != nil {
		return "", err
	}
	return result, nil
}

func writeProfileDataToFile(data bytes.Buffer, dumpType int, dumpPath string) (string, error) {
	file, fileName, err := getBinaryFileNameAndCreate(dumpPath, dumpType)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err = file.Write(data.Bytes()); err != nil {
		return "", err
	}
	return fileName, nil
}

func getBinaryFileNameAndCreate(dump string, dumpType int) (*os.File, string, error) {

	filepath := getBinaryFileName(dump, dumpType)
	_ = os.Remove(filepath)
	f, err := os.OpenFile(filepath, defaultLoggerFlags, defaultLoggerPerm)
	if err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dump, 0o755); err != nil {
			return nil, filepath, err
		}
		f, err = os.OpenFile(filepath, defaultLoggerFlags, defaultLoggerPerm)
		if err != nil {
			return nil, filepath, err
		}
	}
	return f, filepath, err
}

func getBinaryFileName(filePath string, dumpType int) string {
	return path.Join(filePath, check2name[dumpType]+".pprof")
}
