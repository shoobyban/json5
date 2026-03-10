package main

import (
	"fmt"

	"github.com/shoobyban/json5"
)

func main() {
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
		"favorites": ["pizza", 42, false, null, {"item": "book", price: 10.99, "in_stock": true,}],
		        key24: "Hello, \U{0x1F600}world!",
	}`

	// Decode into a generic map using the encoding/json-compatible API
	var result map[string]any
	err := json5.Decode([]byte(input), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result)

	// Encode with indentation using the encoding/json-compatible API
	m, err := json5.EncodeIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(m))

	// Decode an array
	var result2 []any
	err = json5.Decode([]byte(`["a", 'b', 1]`), &result2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result2)

	// Decode a scalar string
	var result3 string
	err = json5.Decode([]byte(`'a'`), &result3)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result3)

	// Decode a scalar number
	var result4 int
	err = json5.Decode([]byte(`42`), &result4)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result4)

	// Validate JSON5
	fmt.Printf("Valid: %v\n", json5.Valid([]byte(`{key: "value",}`)))
	fmt.Printf("Valid: %v\n", json5.Valid([]byte(`{key: }`)))
}
