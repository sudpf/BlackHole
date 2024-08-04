package common

import (
	"reflect"
)

func FieldName(v interface{}, fieldPtr interface{}) string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if fieldPtr is a pointer to a field
	fieldVal := reflect.ValueOf(fieldPtr)
	if fieldVal.Kind() != reflect.Ptr {
		return ""
	}
	fieldVal = fieldVal.Elem()

	// 获取字段所在的结构体类型
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Addr().Interface() == fieldPtr {
			return val.Type().Field(i).Name
		}
	}
	return ""
}
