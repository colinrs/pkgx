package structx

import (
	"errors"
	"reflect"
	"sort"
)

var (
	// ErrUnSupportHaystack ...
	ErrUnSupportHaystack = errors.New("haystack must be slice, array or map")
)

// In ...
func In(haystack interface{}, needle interface{}) (bool, error) {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true, nil
			}
		}

		return false, nil
	}

	return false, ErrUnSupportHaystack
}

// InIntSlice ...
func InIntSlice(haystack []int, needle int) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

// InStringSlice ...
func InStringSlice(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

// InIntSliceSortedFunc ...
func InIntSliceSortedFunc(haystack []int) func(int) bool {
	sort.Ints(haystack)

	return func(needle int) bool {
		index := sort.SearchInts(haystack, needle)
		return index < len(haystack) && haystack[index] == needle
	}
}

// InStringSliceSortedFunc ...
func InStringSliceSortedFunc(haystack []string) func(string) bool {
	sort.Strings(haystack)

	return func(needle string) bool {
		index := sort.SearchStrings(haystack, needle)
		return index < len(haystack) && haystack[index] == needle
	}
}

// SortInIntSlice ...
func SortInIntSlice(haystack []int, needle int) bool {
	sort.Ints(haystack)

	index := sort.SearchInts(haystack, needle)
	return index < len(haystack) && haystack[index] == needle
}

// SortInStringSlice ...
func SortInStringSlice(haystack []string, needle string) bool {
	sort.Strings(haystack)

	index := sort.SearchStrings(haystack, needle)
	return index < len(haystack) && haystack[index] == needle
}

// InIntSliceMapKeyFunc ...
func InIntSliceMapKeyFunc(haystack []int) func(int) bool {
	set := make(map[int]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	return func(needle int) bool {
		_, ok := set[needle]
		return ok
	}
}

// InStringSliceMapKeyFunc ...
func InStringSliceMapKeyFunc(haystack []string) func(string) bool {
	set := make(map[string]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	return func(needle string) bool {
		_, ok := set[needle]
		return ok
	}
}

// MapKeyInIntSlice ...
func MapKeyInIntSlice(haystack []int, needle int) bool {
	set := make(map[int]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	_, ok := set[needle]
	return ok
}

// MapKeyInStringSlice ...
func MapKeyInStringSlice(haystack []string, needle string) bool {
	set := make(map[string]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	_, ok := set[needle]
	return ok
}
