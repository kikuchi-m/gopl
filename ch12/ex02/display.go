package display

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Display(name string, x interface{}) {
	DisplayRecursive(name, x, 10)
}

func DisplayRecursive(name string, x interface{}, maxDepth int) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0, maxDepth)
}

// formatAtom formats a value without inspecting its internal structure.
// It is a copy of the the function in gopl.io/ch11/format.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)

		// ++ ch12.ex01
	case reflect.Struct:
		var s strings.Builder
		s.WriteString(fmt.Sprintf("%s{", v.Type())) // ignore errors
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				s.WriteString(", ")
			}
			s.WriteString(fmt.Sprintf("%s=%s", v.Type().Field(i).Name, formatAtom(v.Field(i))))
		}
		s.WriteString("}")
		return s.String()

	case reflect.Array:
		var s strings.Builder
		s.WriteString(fmt.Sprintf("%s{", v.Type()))
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				s.WriteString(", ")
			}
			s.WriteString(fmt.Sprintf("%s", formatAtom(v.Index(i))))
		}
		s.WriteString("}")
		return s.String()
		// ++ ch12.ex01

	default: // reflect.Interface
		return v.Type().String() + " value"
	}
}

//!+display
func display(path string, v reflect.Value, depth, maxDepth int) {
	if depth > maxDepth {
		fmt.Printf("%s = %s\n", path, formatAtom(v))
		return
	}
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth+1, maxDepth)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), depth+1, maxDepth)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key), depth+1, maxDepth)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), depth+1, maxDepth)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), depth+1, maxDepth)
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}

//!-display
