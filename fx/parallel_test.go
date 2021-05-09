package fx

import (
	"fmt"
	"testing"
)

func TestParallel(t *testing.T) {
	type args struct {
		fns []func()
	}
	type tcase struct {
		name string
		args args
	}
	tc := tcase{
		name:"TestParallel",
		args: args{
			[]func(){
			func(){
			fmt.Printf("1111111")
		}}},
	}
	tests := []tcase {tc}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
