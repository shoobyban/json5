# JSON5 Parser for Go and TinyGo

This is a simple JSON5 parser (and tokenizer if the parser is not good enough) implemented in Go, supporting full functionality in [TinyGo](https://tinygo.org/).

It supports the JSON5 specification, including unquoted keys, escape sequences in strings, booleans, `null`, and hexadecimal numbers.

## Features

As TinyGo does not support reflection, the parser does not use reflection to convert JSON5 tokens into Go native types. Instead, it uses a simple recursive descent parser to convert JSON5 tokens into `map[string]interface{}`, `[]interface{}`, and `string`, `int`, `float64`, `bool`, `nil`.
As I'm lazy, it uses strings as input, feel free to change it to `io.Reader` if you want to parse large files (not an issue for TinyGo or at least in my use case). Pull requests are always welcome.

- **JSON5 Tokenizer**:
  - Handles basic JSON5 syntax: braces, brackets, commas, colons.
  - Supports single-line (`//`) and multi-line (`/* ... */`) comments.
  - Recognizes strings, numbers, booleans (`true`, `false`), and `null` (returns `nil`).
  - Supports unquoted keys in objects.
  - Parses escape sequences in strings, including `\n`, `\t`, `\\`, etc.
  - Parses hexadecimal numbers (e.g., `0x1E`).
  - Parses Unicode escape sequences in strings (e.g., `\u{1F600}`, `\U0X1F4A9`).

## Example

### Parser

```go
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

	result, err := json5.UnMarshal(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Parsed result: %+v\n", result)
}
```

Output:
```
Parsed result: map[address:map[city:New York zipcode:10001] age:42 children:<nil> favorites:[pizza 42 false <nil> map[in_stock:true item:book price:10.99]] hexadecimal:912559 key24:Hello, 😀world! lineBreaks:line 1
line 2
                        line 3 married:true name:John Doe]
```

### Only the tokenizer

```go
package main

import (
	"fmt"

	"github.com/shoobyban/json5"
)

func main() {
	input := `
    {
        // Single-line comment
        "key1": "value1",
        'key2': 'value2',
        key4: 0x1E,
        keyFive: "Hello, \u{0x1F600}world!",
    }`
	tokens := json5.Tokenize(input)
	for _, token := range tokens {
		fmt.Printf("Type: %d, Value: %s\n", token.Type, token.Value)
	}
}
```
Output:
```
Type: 0, Value: {
Type: 11, Value: // Single-line comment
Type: 6, Value: key1
Type: 4, Value: :
Type: 6, Value: value1
Type: 5, Value: ,
Type: 6, Value: key2
Type: 4, Value: :
Type: 6, Value: value2
Type: 5, Value: ,
Type: 6, Value: key4
Type: 4, Value: :
Type: 7, Value: 0x1E
Type: 5, Value: ,
Type: 6, Value: keyFive
Type: 4, Value: :
Type: 6, Value: Hello, 😀world!
Type: 5, Value: ,
Type: 1, Value: }
```

## Benchmarks

Just joking. Horrible performance, but it works for my use case. Feel free to improve it.

```
goos: darwin
goarch: arm64
pkg: github.com/shoobyban/json5
cpu: Apple M1
=== RUN   BenchmarkJSON
BenchmarkJSON
BenchmarkJSON-8           406210              2930 ns/op            5697 B/op         51 allocs/op
PASS
ok      github.com/shoobyban/json5      2.227s
```

## License

This project is licensed under the MIT License.

## Acknowledgments

- [TinyGo](https://tinygo.org/) - you guys rock!
- [JSON5](https://json5.org/)
