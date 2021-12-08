package utils

import (
	"errors"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

var (
	errorNilPointer = errors.New("src_or_dst_cannot_be_nil")
	errorEncoding   = errors.New("encoding_failed_during_deep_copy")
	errorDecoding   = errors.New("decoding_failed_during_deep_copy")
)

// json-iterator/go deep copy
var jsoniterJSON = jsoniter.ConfigFastest

// DeepCopy ...
func DeepCopy(dst interface{}, src interface{}) error {
	return jsoniterDeepCopy(dst, src)
}

func jsoniterDeepCopy(dst interface{}, src interface{}) error {
	if IsPointerPointToNil(dst) || IsPointerPointToNil(src) {
		return errorNilPointer
	}

	encodedBytes, err := jsoniterJSON.Marshal(src)
	if err != nil {
		return errorEncoding
	}
	if err := jsoniterJSON.Unmarshal(encodedBytes, dst); err != nil {
		return errorDecoding
	}

	return nil
}

// IsPointerPointToNil judge whether pointer is point to a nil
func IsPointerPointToNil(i interface{}) bool {
	for {
		if i == nil {
			return true
		}
		if reflect.ValueOf(i).Kind() == reflect.Ptr {
			if reflect.ValueOf(i).IsNil() {
				return true
			}
			i = reflect.ValueOf(i).Elem().Interface()
		} else {
			break
		}
	}
	return false
}
