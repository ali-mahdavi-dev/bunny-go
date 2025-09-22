package helpers

import (
	"reflect"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/helpers/kind"
)

func ToPtr(val reflect.Value) reflect.Value {
	typ := val.Type()
	if !kind.Ptr(typ) {
		// this creates a pointer type inherently
		ptrVal := reflect.New(typ)
		ptrVal.Elem().Set(val)
		val = ptrVal
	}
	return val
}

func FromPtr(val reflect.Value) reflect.Value {
	if kind.Ptr(val.Type()) {
		val = val.Elem()
	}
	return val
}
