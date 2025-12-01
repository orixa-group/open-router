package schema

import (
	"fmt"
	"reflect"
)

type DataType string

const (
	Object  DataType = "object"
	Number  DataType = "number"
	Integer DataType = "integer"
	String  DataType = "string"
	Array   DataType = "array"
	Ptr     DataType = "ptr"
	Boolean DataType = "boolean"
)

func ReflectDataType(t reflect.Type) (DataType, error) {
	kind := t.Kind()
	if dt, ok := dataTypeReflector[kind]; ok {
		return dt, nil
	}

	return "", fmt.Errorf("unsupported type: %s", kind.String())
}

func (d DataType) String() string {
	return string(d)
}

var dataTypeReflector = map[reflect.Kind]DataType{
	reflect.String:  String,
	reflect.Int:     Integer,
	reflect.Int8:    Integer,
	reflect.Int16:   Integer,
	reflect.Int32:   Integer,
	reflect.Int64:   Integer,
	reflect.Uint:    Integer,
	reflect.Uint8:   Integer,
	reflect.Uint16:  Integer,
	reflect.Uint32:  Integer,
	reflect.Uint64:  Integer,
	reflect.Float32: Number,
	reflect.Float64: Number,
	reflect.Bool:    Boolean,
	reflect.Slice:   Array,
	reflect.Array:   Array,
	reflect.Struct:  Object,
	reflect.Ptr:     Ptr,
}
