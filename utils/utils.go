package utils

import "reflect"

// RoundUp to round up float.
func RoundUp(v float64) int {
	if v != float64(int(v)) {
		return int(v) + 1
	}
	return int(v)
}

// isEmptyStruct check value of struct
func IsEmptyStruct(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Struct:
		return field.IsZero()
	default:
		return false
	}
}
