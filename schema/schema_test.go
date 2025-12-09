package schema

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	// ---------------------------------------------------------
	// Helper Structs for Testing
	// ---------------------------------------------------------

	// Simple struct with basic tags
	type SimpleUser struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Struct with omitempty (should affect Required list)
	type OptionalFields struct {
		ID    string `json:"id"`
		Extra string `json:"extra,omitempty"`
	}

	// Nested structs
	type Address struct {
		City string `json:"city"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	// Struct with arrays
	type TaggedPost struct {
		Tags []string `json:"tags"`
	}

	// Struct with ignored fields (no json tag or unexported)
	type PrivateData struct {
		Public  string `json:"public"`
		Ignored string // No tag, should be skipped by objectSchemaReflector
		private string `json:"private"` // Unexported, should be skipped
	}

	// ---------------------------------------------------------
	// Test Cases
	// ---------------------------------------------------------
	tests := []struct {
		name      string
		input     any
		wantType  DataType
		checkFunc func(*testing.T, *Schema)
		wantErr   error
	}{
		{
			name:     "Primitive String",
			input:    "hello",
			wantType: String,
		},
		{
			name:     "Primitive Integer",
			input:    42,
			wantType: Integer,
		},
		{
			name:     "Primitive Boolean",
			input:    true,
			wantType: Boolean,
		},
		{
			name:     "Primitive Float",
			input:    3.14,
			wantType: Number,
		},
		{
			name:     "Slice of Ints",
			input:    []int{1, 2, 3},
			wantType: Array,
			checkFunc: func(t *testing.T, s *Schema) {
				assert.NotNil(t, s.Items, "Array schema must have Items defined")
				assert.Equal(t, Integer, s.Items.Type)
			},
		},
		{
			name:     "Pointer to String",
			input:    new(string),
			wantType: String, // Ptr reflector should dereference to the underlying type
		},
		{
			name:     "Simple Struct",
			input:    SimpleUser{},
			wantType: Object,
			checkFunc: func(t *testing.T, s *Schema) {
				assert.Len(t, s.Properties, 2)
				assert.Contains(t, s.Properties, "name")
				assert.Contains(t, s.Properties, "age")

				// Verify types of properties
				assert.Equal(t, String, s.Properties["name"].Type)
				assert.Equal(t, Integer, s.Properties["age"].Type)

				// Verify Required fields (both should be required as neither has omitempty)
				assert.Contains(t, s.Required, "name")
				assert.Contains(t, s.Required, "age")
			},
		},
		{
			name:     "Struct with omitempty",
			input:    OptionalFields{},
			wantType: Object,
			checkFunc: func(t *testing.T, s *Schema) {
				assert.Contains(t, s.Properties, "id")
				assert.Contains(t, s.Properties, "extra")

				// 'id' is required
				assert.Contains(t, s.Required, "id")
				// 'extra' has omitempty, so it should NOT be in Required
				assert.NotContains(t, s.Required, "extra")
			},
		},
		{
			name:     "Nested Structs",
			input:    Person{},
			wantType: Object,
			checkFunc: func(t *testing.T, s *Schema) {
				addrSchema, ok := s.Properties["address"]
				assert.True(t, ok)
				assert.Equal(t, Object, addrSchema.Type)
				assert.Contains(t, addrSchema.Properties, "city")
			},
		},
		{
			name:     "Struct with Array Field",
			input:    TaggedPost{},
			wantType: Object,
			checkFunc: func(t *testing.T, s *Schema) {
				tags, ok := s.Properties["tags"]
				assert.True(t, ok)
				assert.Equal(t, Array, tags.Type)
				assert.Equal(t, String, tags.Items.Type)
			},
		},
		{
			name:     "Ignored Fields",
			input:    PrivateData{},
			wantType: Object,
			checkFunc: func(t *testing.T, s *Schema) {
				// Only "public" has a valid json tag and is exported
				assert.Len(t, s.Properties, 1)
				assert.Contains(t, s.Properties, "public")

				assert.NotContains(t, s.Properties, "Ignored")
				assert.NotContains(t, s.Properties, "private")
			},
		},
		{
			name:    "Unsupported Map",
			input:   map[string]string{},
			wantErr: errors.New("unsupported type"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.input)

			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.wantType, got.Type)

			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestSchema_MarshalJSON(t *testing.T) {
	// Test that MarshalJSON initializes properties if nil (though encoding/json might still omit empty maps)
	s := &Schema{Type: Object}

	bytes, err := json.Marshal(s)
	assert.NoError(t, err)

	jsonString := string(bytes)
	assert.Contains(t, jsonString, `"type":"object"`)

	// Ensure the custom marshaller doesn't crash on nil properties
	s2 := &Schema{Type: String}
	bytes2, err := json.Marshal(s2)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes2), `"type":"string"`)
}
