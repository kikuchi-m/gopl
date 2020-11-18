package cyclic

import (
	"testing"
	"unsafe"
)

func TestEqual(t *testing.T) {
	one := 1

	type CyclePtr *CyclePtr
	var cyclePtr1 CyclePtr
	cyclePtr1 = &cyclePtr1

	type CycleSlice []CycleSlice
	var cycleSlice = make(CycleSlice, 1)
	cycleSlice[0] = cycleSlice

	type CyclicMap map[int]*CyclicMap
	var cyclicMap = make(CyclicMap)
	cyclicMap[1] = &cyclicMap

	type MyStruct struct {
		s     string
		value *MyStruct
	}
	var myStruct MyStruct
	myStruct.s = "a"
	nonCyclicStruct := MyStruct{"b", &myStruct}
	var cyclicStruct MyStruct
	cyclicStruct.s = "c"
	cyclicStruct.value = &cyclicStruct

	ch1 := make(chan int)
	var ch1ro <-chan int = ch1

	type mystring string

	var iface1 interface{} = &one

	for _, test := range []struct {
		x    interface{}
		want bool
	}{
		// basic types
		{true, false},
		{false, false},
		{1, false},
		{1.2, false},
		{complex(1, 2), false},
		{"foo", false},
		{mystring("foo"), false},
		// // slices
		{[]int{1}, false},
		{[]string{"foo"}, false},
		// slice cycles
		{cycleSlice, true},
		// maps
		{map[string][]int{"foo": {1, 2, 3}}, false},
		{cyclicMap, true},
		// pointers
		{&one, false},
		// pointer cycles
		{cyclePtr1, true},
		// struct
		{nonCyclicStruct, false},
		{cyclicStruct, true},
		// functions
		{(func())(nil), false},
		// arrays
		{[...]int{1, 2, 3}, false},
		// channels
		{ch1, false},
		{ch1ro, false},
		// interfaces
		{&iface1, false},
		// unsafe pointer
		{unsafe.Pointer(&one), false},
		{unsafe.Pointer(cyclePtr1), false}, // ???
	} {
		if Cyclic(test.x) != test.want {
			t.Errorf("Cyclic(%s) = %t", test.x, !test.want)
		}
	}
}
