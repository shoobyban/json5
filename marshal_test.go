package json5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalNil(t *testing.T) {
	result, err := Marshal(nil, "")
	assert.NoError(t, err)
	assert.Equal(t, "null", result)
}

func TestMarshalBool(t *testing.T) {
	result, err := Marshal(true, "")
	assert.NoError(t, err)
	assert.Equal(t, "true", result)

	result, err = Marshal(false, "")
	assert.NoError(t, err)
	assert.Equal(t, "false", result)
}

func TestMarshalNumber(t *testing.T) {
	result, err := Marshal(42, "")
	assert.NoError(t, err)
	assert.Equal(t, "42", result)

	result, err = Marshal(3.14, "")
	assert.NoError(t, err)
	assert.Equal(t, "3.14", result)
}

func TestMarshalString(t *testing.T) {
	result, err := Marshal("Hello\nWorld", "")
	assert.NoError(t, err)
	assert.Equal(t, "\"Hello\\nWorld\"", result)
}

func TestMarshalArray(t *testing.T) {
	input := []interface{}{"a", 1, true, nil}
	expected := `[
"a", 
1, 
true, 
null
]`
	result, err := Marshal(input, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalObject(t *testing.T) {
	input := map[string]interface{}{
		"name": "John Doe",
		"age":  42,
		"address": map[string]interface{}{
			"city":    "New York",
			"zipcode": 10001,
		},
	}
	expected := `{
name: "John Doe",
age: 42,
address: {
city: "New York",
zipcode: 10001,
},
}`
	result, err := Marshal(input, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalComplexObject(t *testing.T) {
	input := map[string]interface{}{
		"simpleKey": "value",
		"complex key": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}
	expected := `{
simpleKey: "value",
"complex key": {
nestedKey: "nestedValue",
},
}`
	result, err := Marshal(input, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
