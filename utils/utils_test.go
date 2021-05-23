package utils

import (
	"reflect"
	"testing"
)

func TestStruct2Map(t *testing.T) {
	type args struct {
		obj interface{}
	}
	type tcase struct{
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}
	type S struct{
		TN string
	}
	tc1:= tcase{
		name:"TestStruct2Map",
		args:args{
			obj: S{
				TN:"xxxxx",
			},
		},
		want:map[string]interface{}{"TN":"xxxxx"},
		wantErr: false,
	}
	tests := []tcase {tc1}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := Struct2Map(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct2Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Struct2Map() got = %v, want %v", got, tt.want)
			}
		})
	}
}
