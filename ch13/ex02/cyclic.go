package cyclic

import (
	"reflect"
	"unsafe"
)

type ptr struct {
	x unsafe.Pointer
	t reflect.Type
}

func cyclic(x reflect.Value, seen map[ptr]bool) bool {
	if !x.IsValid() {
		return false
	}

	if x.CanAddr() {
		xptr := ptr{unsafe.Pointer(x.UnsafeAddr()), x.Type()}
		if seen[xptr] {
			return true
		}
		seen[xptr] = true
	}

	switch x.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.UnsafePointer: // ???
		return false

	case reflect.Ptr, reflect.Interface:
		return cyclic(x.Elem(), seen)

	case reflect.Array, reflect.Slice:
		for i := 0; i < x.Len(); i++ {
			if cyclic(x.Index(i), seen) {
				return true
			}
		}
		return false

	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if cyclic(x.Field(i), seen) {
				return true
			}
		}
		return false

	case reflect.Map:
		for _, k := range x.MapKeys() {
			if cyclic(x.MapIndex(k), seen) {
				return true
			}
		}
		return false

	case reflect.Chan, reflect.Func:
		return false
	}
	panic("unreachable")
}

func Cyclic(x interface{}) bool {
	seen := make(map[ptr]bool)
	return cyclic(reflect.ValueOf(x), seen)
}
