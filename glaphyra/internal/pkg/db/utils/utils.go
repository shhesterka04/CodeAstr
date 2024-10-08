package utils

import "reflect"

func IsZeroVal(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func SetIfNotZero(valMap map[string]interface{}, name string, value interface{}) map[string]interface{} {
	if valMap == nil {
		return valMap
	}

	if !IsZeroVal(value) {
		valMap[name] = value
	}

	return valMap
}
