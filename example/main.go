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

	result, err := json5.Parse(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result)

	result2, err := json5.Parse(`["a", 'b', 1]`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result2)

	result3, err := json5.Parse(`'a'`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result3)

	result4, err := json5.Parse(`42`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result4)

	result5, err := json5.Parse(``)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result5)
}
