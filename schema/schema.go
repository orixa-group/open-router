package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Schema struct {
	Type                 DataType           `json:"type,omitempty"`
	Description          string             `json:"description,omitempty"`
	Enum                 []string           `json:"enum,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	AdditionalProperties any                `json:"additionalProperties,omitempty"`
	Nullable             bool               `json:"nullable,omitempty"`
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	if s.Properties == nil {
		s.Properties = make(map[string]*Schema)
	}
	type Alias Schema
	return json.Marshal(struct {
		Alias
	}{
		Alias: (Alias)(*s),
	})
}

func Generate(v any) (*Schema, error) {
	return reflectSchema(reflect.TypeOf(v))
}

func reflectSchema(t reflect.Type) (*Schema, error) {
	dt, err := ReflectDataType(t)
	if err != nil {
		return nil, err
	}

	for _, r := range schemaReflectors {
		if r.DataType() == dt {
			return r.Schema(t)
		}
	}

	return nil, fmt.Errorf("unsupported type: %s", t.Kind().String())
}

type SchemaReflector interface {
	DataType() DataType
	Schema(t reflect.Type) (*Schema, error)
}

var schemaReflectors = []SchemaReflector{
	stringSchemaReflector{},
	integerSchemaReflector{},
	numberSchemaReflector{},
	boolSchemaReflector{},
	arraySchemaReflector{},
	ptrSchemaReflector{},
	objectSchemaReflector{},
}

type stringSchemaReflector struct{}

func (r stringSchemaReflector) DataType() DataType {
	return String
}

func (r stringSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	return &Schema{Type: String}, nil
}

type integerSchemaReflector struct{}

func (r integerSchemaReflector) DataType() DataType {
	return Integer
}

func (r integerSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	return &Schema{Type: Integer}, nil
}

type numberSchemaReflector struct{}

func (r numberSchemaReflector) DataType() DataType {
	return Number
}

func (r numberSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	return &Schema{Type: Number}, nil
}

type boolSchemaReflector struct{}

func (r boolSchemaReflector) DataType() DataType {
	return Boolean
}

func (r boolSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	return &Schema{Type: Boolean}, nil
}

type arraySchemaReflector struct{}

func (r arraySchemaReflector) DataType() DataType {
	return Array
}

func (r arraySchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	items, err := reflectSchema(t.Elem())
	if err != nil {
		return nil, err
	}

	return &Schema{Type: Array, Items: items}, nil
}

type ptrSchemaReflector struct{}

func (r ptrSchemaReflector) DataType() DataType {
	return Ptr
}

func (r ptrSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	return reflectSchema(t.Elem())
}

type objectSchemaReflector struct{}

func (r objectSchemaReflector) DataType() DataType {
	return Object
}

func (r objectSchemaReflector) Schema(t reflect.Type) (*Schema, error) {
	required := make([]string, 0)
	props := make(map[string]*Schema)

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if tag := field.Tag.Get("json"); tag != "" {
			name := strings.TrimSuffix(tag, ",omitempty")
			if !strings.HasSuffix(tag, ",omitempty") {
				required = append(required, name)
			}

			item, err := reflectSchema(field.Type)
			if err != nil {
				return nil, err
			}

			props[name] = item
		}
	}

	return &Schema{
		Type:                 Object,
		Properties:           props,
		Required:             required,
		AdditionalProperties: false,
	}, nil
}
