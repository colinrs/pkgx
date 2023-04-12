package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStruct2Map(t *testing.T) {
	type args struct {
		obj interface{}
	}
	type tcase struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}
	type S struct {
		TN string
	}
	tc1 := tcase{
		name: "TestStruct2Map",
		args: args{
			obj: S{
				TN: "xxxxx",
			},
		},
		want:    map[string]interface{}{"TN": "xxxxx"},
		wantErr: false,
	}
	tests := []tcase{tc1}
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

type Person struct {
	Name string
}

func TestSetValueV2(t *testing.T) {
	var sInt *int
	err := SetValue(2, &sInt)
	if err != nil {
		return
	}
	if sInt != nil && *sInt == 2 {
		fmt.Printf("s:%d\n", *sInt)
	}

	var sBool *bool
	err = SetValue(true, &sBool)
	if err != nil {
		return
	}
	if sBool != nil {
		fmt.Printf("s:%v\n", *sBool)
	}

	var sStruct1 *Person
	err = SetValue(&Person{
		Name: "sStruct1",
	}, &sStruct1)
	if err != nil {
		return
	}
	if sStruct1 != nil {
		fmt.Printf("sStruct11:%+v\n", sStruct1)
	}

	var sStruct2 *Person
	err = SetValue(Person{
		Name: "sStruct2",
	}, &sStruct2)
	if err != nil {
		return
	}
	if sStruct2 != nil {
		fmt.Printf("sStruct2:%+v\n", sStruct2)
	}

}
