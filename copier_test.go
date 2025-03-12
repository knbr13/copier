package main

import (
	"reflect"
	"testing"
	"unsafe"
)

type testStruct struct {
	Int    int
	Float  float64
	Str    string
	Arr    [3]int
	Slice  []int
	Map    map[string]int
	Ptr    *int
	Chan   chan int
	Struct *testStruct
}

func TestShallowCopyStruct(t *testing.T) {
	i := 14
	c := make(chan int, 2)
	c <- 14
	m := map[string]int{"a": 1, "b": 2}
	s := []int{4, 5, 6}
	ts := testStruct{
		Int: 10,
	}
	tests := []struct {
		name     string
		src      testStruct
		expected testStruct
	}{
		{
			name: "shallow copy struct",
			src: testStruct{
				Int:    1,
				Float:  2.3,
				Str:    "hello",
				Arr:    [3]int{1, 2, 3},
				Slice:  s,
				Map:    m,
				Ptr:    &i,
				Chan:   c,
				Struct: &ts,
			},
			expected: testStruct{
				Int:    1,
				Float:  2.3,
				Str:    "hello",
				Arr:    [3]int{1, 2, 3},
				Slice:  s,
				Map:    m,
				Ptr:    &i,
				Chan:   c,
				Struct: &ts,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dst := &testStruct{}
			err := ShallowCopyStruct(dst, test.src)
			if err != nil {
				t.Errorf("ShallowCopyStruct() error = %v", err)
				return
			}

			if test.expected.Arr != dst.Arr {
				t.Errorf("dst.Arr = %v, want %v", dst.Arr, test.expected.Arr)
			}
			if !sameMemoryAddress(test.expected.Slice, dst.Slice) || len(test.expected.Slice) != len(dst.Slice) {
				t.Errorf("dst.Slice = %v, want %v", dst.Slice, test.expected.Slice)
			}
			if !sameMemoryAddress(test.expected.Map, dst.Map) || len(test.expected.Map) != len(dst.Map) {
				t.Errorf("dst.Map = %v, want %v", dst.Map, test.expected.Map)
			}
			if test.expected.Ptr != dst.Ptr {
				t.Errorf("dst.Ptr = %v, want %v", dst.Ptr, test.expected.Ptr)
			}
			if test.expected.Chan != dst.Chan {
				t.Errorf("dst.Chan = %v, want %v", dst.Chan, test.expected.Chan)
			}
			if test.expected.Struct != dst.Struct {
				t.Errorf("dst.Struct = %v, want %v", dst.Struct, test.expected.Struct)
			}
		})
	}
}

func TestDeepCopyStruct(t *testing.T) {
	i, _ := 13, 13
	c := make(chan int, 2)
	var t1, t2 testStruct
	t1.Int = 10
	t2.Int = 10
	tests := []struct {
		name     string
		src      testStruct
		expected testStruct
	}{
		{
			name: "deep copy struct",
			src: testStruct{
				Int:   1,
				Float: 2.3,
				Str:   "hello",
				Arr:   [3]int{1, 2, 3},
				Slice: []int{4, 5, 6},
				Map:   map[string]int{"a": 1, "b": 2},
				Ptr:   &i,
				Chan:  c,
				// Chan:   make(chan int),
				Struct: &t1,
			},
			expected: testStruct{
				Int:   1,
				Float: 2.3,
				Str:   "hello",
				Arr:   [3]int{1, 2, 3},
				Slice: []int{4, 5, 6},
				Map:   map[string]int{"a": 1, "b": 2},
				Ptr:   &i,
				Chan:  c,
				// Chan:   make(chan int),
				Struct: &t1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dst := &testStruct{}
			err := DeepCopyStruct(dst, test.src)
			if err != nil {
				t.Errorf("DeepCopyStruct() error = %v", err)
				return
			}

			if !reflect.DeepEqual(*dst, test.expected) {
				t.Errorf("DeepCopyStruct() \ngot = %v, \nexpected = %v", dst, test.expected)
			}

			// Modify the original struct and check if the copied struct remains unchanged
			// test.src.Int = 100
			// test.src.Slice[0] = 100
			// test.src.Map["a"] = 100
			// *test.src.Ptr = 100
			// test.src.Struct.Int = 100

			// if reflect.DeepEqual(dst, test.expected) {
			// 	t.Errorf("DeepCopyStruct() dst was modified after modifying src")
			// }
		})
	}
}

func sameMemoryAddress(slice1, slice2 interface{}) bool {
	ptr1 := unsafe.Pointer(reflect.ValueOf(slice1).Pointer())
	ptr2 := unsafe.Pointer(reflect.ValueOf(slice2).Pointer())
	return ptr1 == ptr2
}
