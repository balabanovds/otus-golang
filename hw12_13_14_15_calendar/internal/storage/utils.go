package storage

import "reflect"

func IsZeroValue(v interface{}) bool {
	if v == nil {
		return false
	}

	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
