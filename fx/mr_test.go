package fx

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func testFuncForNoneArgs() {
	fmt.Print("t0 name:T0\n")
}

func testFuncForOneArgs(name string) {
	fmt.Printf("t1 name:%s\n", name)
}

func testFuncForOneArgsAndErr(name string) error {
	fmt.Printf("t2 name:%s\n", name)
	return nil
}

func testFuncForOneArgsAndResultWithErr(name string) (string, error) {
	fmt.Printf("t3 name:%s\n", name)
	return fmt.Sprintf("%s_%s", "T3", name), nil
}

func testFuncForThreeArgsAndResultWithErr(name string, name2 string, sleep time.Duration) (string, error) {
	time.Sleep(sleep)
	fmt.Printf("t4 name:%s, name2:%s\n", name, name2)
	return fmt.Sprintf("%s_%s_%s", "T4", name, name2), errors.New("xxx")
}

func TestConcurrent_Run(t *testing.T) {
	type tcast struct {
		name       string
		Concurrent *Concurrent
	}
	concurrentSuccess := NewConcurrent(WithTimeout(10*time.Second), WithGoConcurrentLimit(2))
	name := "TestForSuccess"
	concurrentSuccess.AddFunc(testFuncForNoneArgs, nil)
	concurrentSuccess.AddFunc(testFuncForOneArgs, concurrentSuccess.GetParamValues(name))
	concurrentSuccess.AddFunc(testFuncForOneArgsAndErr, concurrentSuccess.GetParamValues(name))
	concurrentSuccess.AddFunc(testFuncForOneArgsAndResultWithErr, concurrentSuccess.GetParamValues(name))
	concurrentSuccess.AddFunc(testFuncForThreeArgsAndResultWithErr, concurrentSuccess.GetParamValues(name, name, 5*time.Second))
	tcSuccess := tcast{
		name:       name,
		Concurrent: concurrentSuccess,
	}
	name = "TestForFailed"
	concurrentFailed := NewConcurrent(WithTimeout(3 * time.Second))
	concurrentFailed.AddFunc(testFuncForNoneArgs, nil)
	concurrentFailed.AddFunc(testFuncForOneArgs, concurrentFailed.GetParamValues(name))
	concurrentFailed.AddFunc(testFuncForOneArgsAndErr, concurrentFailed.GetParamValues(name))
	concurrentFailed.AddFunc(testFuncForOneArgsAndResultWithErr, concurrentFailed.GetParamValues(name))
	concurrentFailed.AddFunc(testFuncForThreeArgsAndResultWithErr, concurrentFailed.GetParamValues(name, name, 5*time.Second))
	tcFailed := tcast{
		name:       name,
		Concurrent: concurrentFailed,
	}
	tests := []tcast{tcSuccess, tcFailed}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Concurrent.Run()
			result, err := tt.Concurrent.Result()
			if err != nil {
				fmt.Printf("Concurrent err:%s\n", err.Error())
				return
			}
			for index, values := range result {
				fmt.Printf("============== index:%d\n", index)
				for _, value := range values {
					fmt.Printf("Result:%T,%+v\n", value, value)
				}
			}
			return
		})
	}
}
