package json5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSON5Object(t *testing.T) {
	input := `{
        "name": "John Doe", 
        "age": 42,
        "married": true, 
        "children": null, 
          hexadecimal: 0xdecaf,
        "address": {
            "city": "New York",
            "zipcode": 10001
        },
        lineBreaks: "line 1\nline 2
		line 3",
		unicode: "Hello beauty! -\U{0x1F600}-",
        "favorites": ["pizza", 42, false, null, {"item": "book", price: 10.99, "in_stock": true,}]
    }`

	result, err := UnMarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 object: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}

	expected := map[string]interface{}{
		"address":   map[string]interface{}{"city": "New York", "zipcode": 10001},
		"age":       42,
		"children":  nil,
		"favorites": []interface{}{"pizza", 42, false, nil, map[string]interface{}{"in_stock": true, "item": "book", "price": 10.99}},
		"married":   true,
		"name":      "John Doe", "hexadecimal": 912559, "lineBreaks": "line 1\nline 2\n\t\tline 3",
		"unicode": "Hello beauty! -😀-"}
	assert.Equal(t, expected, result)
}

func TestParseJSON5Array(t *testing.T) {
	input := `["a", 'b', 1]`

	result, err := UnMarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 array: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func TestParseJSON5String(t *testing.T) {
	input := `'a'`

	result, err := UnMarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 string: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func TestParseJSON5Number(t *testing.T) {
	input := `42`

	result, err := UnMarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 number: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func BenchmarkJSON(b *testing.B) {
	input := `{
			"name": "John Doe", 
			"age": 42,
			"married": true, 
			"children": null, 
			  hexadecimal: 0xdecaf,
			"address": {
				"city": "New York",
				"zipcode": 10001
			},
			lineBreaks: "line 1\nline 2
			line 3",
			unicode: "Hello beauty! -\U{0x1F600}-",
			"favorites": ["pizza", 42, false, null, {"item": "book", price: 10.99, "in_stock": true,}]
		}`

	for i := 0; i < b.N; i++ {
		UnMarshal(input)
	}
}
