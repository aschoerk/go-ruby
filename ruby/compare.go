package ruby

import (
	"reflect"
)

// CompareGenerally compares two multidimensional slices of ordered types
// and returns true if all corresponding elements are equal or less than.
func CompareGenerally[T any](a T, b T) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	return compareReflectValues(va, vb)
}

func compareReflectValues(a, b reflect.Value) bool {
	// Handle nil pointers
	if !a.IsValid() || !b.IsValid() {
		return !b.IsValid() // a <= b if b is also nil
	}
	if a.Kind() == reflect.Ptr && a.IsNil() {
		return true // nil is considered the lowest value
	}
	if b.Kind() == reflect.Ptr && b.IsNil() {
		return false // a > nil (b)
	}

	// Dereference pointers
	if a.Kind() == reflect.Ptr {
		a = a.Elem()
	}
	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}

	if a.Kind() != b.Kind() {
		return false
	}

	switch a.Kind() {
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return a.Len() <= b.Len()
		}
		for i := 0; i < a.Len(); i++ {
			if !compareReflectValues(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() <= b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() <= b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() <= b.Float()
	case reflect.String:
		return a.String() <= b.String()
	default:
		return false // Unsupported type
	}
}

// CountDimensions returns the number of dimensions in a multidimensional slice
func CountDimensions(v interface{}) int {
	return countDimensionsRecursive(reflect.ValueOf(v))
}

func countDimensionsRecursive(v reflect.Value) int {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 1 // Consider nil pointer as 1 dimension
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return 0
	}
	if v.Len() == 0 {
		return 1 // Empty slice/array, assume 1 dimension
	}
	return 1 + countDimensionsRecursive(v.Index(0))
}

func isComparable(v interface{}) bool {
	t := reflect.TypeOf(v)

	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.Chan, reflect.UnsafePointer:
		return true
	case reflect.Ptr, reflect.Interface:
		// Pointers and interfaces are comparable if their underlying types are comparable.
		return isComparable(reflect.Indirect(reflect.ValueOf(v)).Interface())
	case reflect.Array:
		// Arrays are comparable if their element types are comparable.
		return isComparable(reflect.New(t.Elem()).Elem().Interface())
	case reflect.Struct:
		// Structs are comparable if all their fields are comparable.
		for i := 0; i < t.NumField(); i++ {
			if !isComparable(reflect.New(t.Field(i).Type).Elem().Interface()) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
