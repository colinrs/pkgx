package fx

import (
	"fmt"
	"testing"
)

func T0() {
	fmt.Print("t1 name:T0\n")
}

func T1(name string) {
	fmt.Printf("t1 name:%s\n", name)
}

func T2(name string) error {
	fmt.Printf("t2 name:%s\n", name)
	return nil
}

func T3(name string) (string, error) {
	fmt.Printf("t3 name:%s\n", name)
	return name, nil
}

func TestMR_Run(t *testing.T) {
	type tcast struct {
		name string
		MR   *MR
	}
	mr1 := NewMR()
	name := "name"
	mr1.AddFunc(T0)
	mr1.AddFunc(T1, []interface{}{name})
	mr1.AddFunc(T2, []interface{}{name})
	mr1.AddFunc(T3, []interface{}{name})
	tc := tcast{
		name: "TestFor1",
		MR:   mr1,
	}
	tests := []tcast{tc}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MR.Run()
			fmt.Printf("Result:%+v\n", tt.MR.Result())
		})
	}
}
