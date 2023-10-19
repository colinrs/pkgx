package structx

// Set ... TODO https://github.com/deckarep/golang-set/blob/master/set.go
type Set interface {
	Add(i interface{}) bool

	Clear()

	Clone() Set

	Contains(i ...interface{}) bool

	Difference(other Set) Set

	Equal(other Set) bool

	Each(func(interface{}) bool)

	Remove(i interface{})

	String() string

	Union(other Set) Set

	Pop() interface{}

	ToSlice() []interface{}
}
