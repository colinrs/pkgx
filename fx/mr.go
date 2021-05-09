package fx

import (
	"fmt"
	"reflect"
)

type Func func()

type MR struct {
	waitGroup RoutineGroup
	fgsCount  uint64
	fgs []Func
	result []interface{}
	run bool
}

func NewMR() *MR {
	return &MR{
		waitGroup: RoutineGroup{},
		fgsCount:  0,
	}
}

func (mr *MR) AddFunc(impl interface{}, args ...[]interface{}) {
	index := mr.fgsCount
	ff := func(){
		handlerFuncType := reflect.TypeOf(impl)
		if handlerFuncType.Kind()!=reflect.Func{
			panic("not a func")
		}
		f := reflect.ValueOf(impl)
		var realParam []reflect.Value
		for _, a := range args {
			realParam = append(realParam, reflect.ValueOf(a))
		}
		mr.result[index] = f.Call(realParam)
	}
	mr.fgsCount++
	mr.fgs = append(mr.fgs, ff)
}

func (mr *MR) Run() {
	mr.run = true
	mr.result = make([]interface{}, len(mr.fgs))
	for _, fn := range mr.fgs{
		mr.waitGroup.RunSafe(fn)
	}
	mr.waitGroup.Wait()
}

func (mr *MR) Result() []interface{}{
	if !mr.run{
		fmt.Print("no result, run it first\n")
		return []interface{}{}

	}
	return mr.result
}