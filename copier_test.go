package copier

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type Basic struct {
	Name    string
	Age     int
	Enabled bool
	Score   float64
}

type PtrStruct struct {
	IntPtr    *int
	StringPtr *string
}

type CollectionStruct struct {
	Slice []int
	Map   map[string]int
}

type Inner struct {
	Value int
}

type Outer struct {
	Inner     Inner
	InnerPtr  *Inner
	InnerList []Inner
}

type Private struct {
	public  string
	Public  string
	private int
}

type TimeStruct struct {
	T time.Time
}

type InterfaceStruct struct {
	Reader interface{}
}

type Embedded struct {
	Basic
	Extra string
}

func TestShallowCopyStruct(t *testing.T) {
	num := 42
	str := "hello"
	now := time.Now()

	tests := []struct {
		name    string
		src     interface{}
		dst     interface{}
		want    interface{}
		wantErr bool
		modify  func(src interface{})
	}{
		{
			name: "basic struct",
			src:  &Basic{Name: "Alice", Age: 30},
			dst:  &Basic{},
			want: &Basic{Name: "Alice", Age: 30},
		},
		{
			name: "pointer fields",
			src:  &PtrStruct{IntPtr: &num, StringPtr: &str},
			dst:  &PtrStruct{},
			want: &PtrStruct{IntPtr: &num, StringPtr: &str},
		},
		{
			name: "slice and map sharing",
			src: &CollectionStruct{
				Slice: []int{1, 2, 3},
				Map:   map[string]int{"a": 1},
			},
			dst: &CollectionStruct{},
			modify: func(src interface{}) {
				s := src.(*CollectionStruct)
				s.Slice[0] = 99
				s.Map["a"] = 99
			},
			want: &CollectionStruct{
				Slice: []int{99, 2, 3},
				Map:   map[string]int{"a": 99},
			},
		},
		{
			name: "time struct",
			src:  &TimeStruct{T: now},
			dst:  &TimeStruct{},
			want: &TimeStruct{T: now},
		},
		{
			name:    "nil destination",
			src:     &Basic{},
			dst:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ShallowCopyStruct(tt.dst, tt.src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ShallowCopyStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.modify != nil {
				tt.modify(tt.src)
			}

			if !tt.wantErr {
				assertShallowCopy(t, tt.dst, tt.want)
			}
		})
	}
}

func TestDeepCopyStruct(t *testing.T) {
	num := 42
	str := "hello"
	testErr := errors.New("test error")

	tests := []struct {
		name    string
		src     interface{}
		dst     interface{}
		want    interface{}
		wantErr bool
		modify  func(src interface{})
	}{
		{
			name: "basic struct",
			src:  &Basic{Name: "Bob", Age: 40},
			dst:  &Basic{},
			want: &Basic{Name: "Bob", Age: 40},
		},
		{
			name: "pointer fields",
			src:  &PtrStruct{IntPtr: &num, StringPtr: &str},
			dst:  &PtrStruct{},
			want: &PtrStruct{IntPtr: &num, StringPtr: &str},
		},
		{
			name: "nested structures",
			src: &Outer{
				Inner:     Inner{Value: 10},
				InnerPtr:  &Inner{Value: 20},
				InnerList: []Inner{{Value: 30}},
			},
			dst: &Outer{},
			modify: func(src interface{}) {
				s := src.(*Outer)
				s.Inner.Value = 99
				s.InnerPtr.Value = 99
				s.InnerList[0].Value = 99
			},
			want: &Outer{
				Inner:     Inner{Value: 10},
				InnerPtr:  &Inner{Value: 20},
				InnerList: []Inner{{Value: 30}},
			},
		},
		{
			name: "unexported fields",
			src:  &Private{public: "secret", Public: "open"},
			dst:  &Private{},
			want: &Private{Public: "open"},
		},
		{
			name: "interface field",
			src:  &InterfaceStruct{Reader: testErr},
			dst:  &InterfaceStruct{},
			want: &InterfaceStruct{Reader: testErr},
		},
		{
			name:    "invalid source type",
			src:     "not a struct",
			dst:     &Basic{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeepCopyStruct(tt.dst, tt.src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeepCopyStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.modify != nil {
				tt.modify(tt.src)
			}

			if !tt.wantErr {
				assertDeepCopy(t, tt.dst, tt.want)
			}
		})
	}
}

func assertShallowCopy(t *testing.T, dst, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("ShallowCopy result mismatch:\nGot: %+v\nWant: %+v", dst, want)
	}
}

func assertDeepCopy(t *testing.T, dst, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(dst, want) {
		t.Errorf("DeepCopy result mismatch:\nGot: %+v\nWant: %+v", dst, want)
	}

	dstVal := reflect.ValueOf(dst).Elem()
	wantVal := reflect.ValueOf(want).Elem()
	checkDeepCopyPointers(t, dstVal, wantVal)
}

func checkDeepCopyPointers(t *testing.T, dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Ptr:
		if !src.IsNil() && dst.Pointer() == src.Pointer() {
			t.Errorf("Found shared pointer at %s", src.Type())
		}
		if !src.IsNil() {
			checkDeepCopyPointers(t, dst.Elem(), src.Elem())
		}

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			checkDeepCopyPointers(t, dst.Field(i), src.Field(i))
		}

	case reflect.Slice:
		if src.Len() > 0 && dst.Pointer() == src.Pointer() {
			t.Error("Slice shares underlying array")
		}
		for i := 0; i < src.Len(); i++ {
			checkDeepCopyPointers(t, dst.Index(i), src.Index(i))
		}

	case reflect.Map:
		if src.Len() > 0 && dst.Pointer() == src.Pointer() {
			t.Error("Map shares underlying data")
		}
	}
}
