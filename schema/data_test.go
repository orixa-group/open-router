package schema

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectDataType(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    DataType
		wantErr error
	}{
		{"int", int(0), Integer, nil},
		{"int8", int8(0), Integer, nil},
		{"int16", int16(0), Integer, nil},
		{"int32", int32(0), Integer, nil},
		{"int64", int64(0), Integer, nil},
		{"uint", uint(0), Integer, nil},
		{"uint8", uint8(0), Integer, nil},
		{"uint16", uint16(0), Integer, nil},
		{"uint32", uint32(0), Integer, nil},
		{"uint64", uint64(0), Integer, nil},
		{"float32", float32(0), Number, nil},
		{"float64", float64(0), Number, nil},
		{"bool", true, Boolean, nil},
		{"string", "test", String, nil},
		{"slice", []int{}, Array, nil},
		{"array", [3]int{}, Array, nil},
		{"struct", struct{ Name string }{}, Object, nil},
		{"pointer", (*int)(nil), Ptr, nil},
		{"map", map[string]int{}, "", errors.New("unsupported type")},
		{"channel", make(chan int), "", errors.New("unsupported type")},
		{"func", func() {}, "", errors.New("unsupported type")},
		{"complex64", complex64(0), "", errors.New("unsupported type")},
		{"complex128", complex128(0), "", errors.New("unsupported type")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReflectDataType(reflect.TypeOf(tt.input))

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
