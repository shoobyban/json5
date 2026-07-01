# JSON5 Parser for Go and TinyGo

[![Go Reference](https://pkg.go.dev/badge/github.com/shoobyban/json5.svg)](https://pkg.go.dev/github.com/shoobyban/json5)

This is a simple JSON5 parser (and tokenizer if the parser is not good enough) implemented in Go, supporting full functionality in [TinyGo](https://tinygo.org/).

It supports the JSON5 specification, including unquoted keys, escape sequences in strings, booleans, `null`, and hexadecimal numbers.

Development currently targets Go 1.25.x. The TinyGo example Makefile defaults to `go1.25.4` for reproducible local builds.

## Features

As TinyGo does not support reflection, the parser does not use reflection to convert JSON5 tokens into Go native types. Instead, it uses a simple recursive descent parser to convert JSON5 tokens into `map[string]any`, `[]any`, and `string`, `int`, `float64`, `bool`, `nil`.
As I'm lazy, it uses strings as input, feel free to change it to `io.Reader` if you want to parse large files (not an issue for TinyGo or at least in my use case). Pull requests are always welcome.

- **JSON5 Tokenizer**:
  - Handles basic JSON5 syntax: braces, brackets, commas, colons.
  - Supports single-line (`//`) and multi-line (`/* ... */`) comments.
  - Recognizes strings, numbers, booleans (`true`, `false`), and `null` (returns `nil`).
  - Supports unquoted keys in objects, including keyword-style keys such as `true`, `false`, `null`, `Infinity`, and `NaN`.
  - Parses escape sequences in strings, including `\n`, `\t`, `\\`, etc.
  - Parses hexadecimal numbers (e.g., `0x1E`).
  - Parses Unicode escape sequences in strings (e.g., `\u{1F600}`, `\U0X1F4A9`).
  - Rejects empty or comment-only documents and unterminated quoted strings as invalid input.

## encoding/json-Compatible API

For users migrating from `encoding/json` or who prefer its conventions, the library provides a compatibility layer with familiar signatures. These functions accept `[]byte` and decode into typed Go values, including structs with `json` tags.

### Decode (like json.Unmarshal)

```go
type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

var cfg Config
err := json5.Decode([]byte(`{
    // server settings
    host: "localhost",
    port: 8080,
}`), &cfg)
// cfg.Host == "localhost", cfg.Port == 8080
```

`Decode` supports `*any`, `*map[string]any`, `*[]any`, scalar pointers (`*string`, `*bool`, `*int`, `*float64`), and structs (via `encoding/json` round-trip).

### Encode / EncodeIndent (like json.Marshal / json.MarshalIndent)

```go
data := map[string]any{"name": "Alice", "active": true}

b, err := json5.Encode(data)
// b == []byte(`{active: true, name: "Alice"}`)

b, err = json5.EncodeIndent(data, "", "  ")
// b == []byte("{\n  active: true,\n  name: \"Alice\",\n}")
```

### Valid (like json.Valid)

```go
json5.Valid([]byte(`{key: "value",}`))  // true  (trailing comma OK)
json5.Valid([]byte(`{key: }`))          // false
json5.Valid([]byte(``))                 // false
json5.Valid([]byte(`// comment only`))  // false
```

### API Comparison

| `encoding/json`                             | `json5` equivalent       |
| ------------------------------------------- | ------------------------ |
| `json.Unmarshal([]byte, any) error`         | `json5.Decode`           |
| `json.Marshal(any) ([]byte, error)`         | `json5.Encode`           |
| `json.MarshalIndent(any, "", " ") (…)`      | `json5.EncodeIndent`     |
| `json.Valid([]byte) bool`                    | `json5.Valid`            |
| `json.Unmarshal` into `any`                 | `json5.Unmarshal` (original) |
| `json.Marshal` returning `string`           | `json5.Marshal` (original)   |

## Deprecation Notice

`UnMarshal` (capital M) is **deprecated**. Use `Unmarshal` or `Decode` instead. `UnMarshal` will continue to work but may be removed in a future major version.

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

	result, err := json5.Unmarshal(input)
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

### Marshalling

```go
package main

import (
    "fmt"

    "github.com/shoobyban/json5"
)

func main() {
    data := map[string]any{
        "name": "John Doe",
        "age": 42,
        "married": true,
        "children": nil,
        "address": map[string]any{
            "city": "New York",
            "zipcode": 10001,
        },
        "favorites": []any{
            "pizza", 42, false, nil,
            map[string]any{
                "item": "book",
                "price": 10.99,
                "in_stock": true,
            },
        },
    }

    // there is a Marshal(value any) function as well that does not indent the output
    result, err := json5.MarshalIndent(data, "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(result)
}
```

Output:
```
{
  name: "John Doe",
  age: 42,
  married: true,
  children: null,
  address: {
    city: "New York",
    zipcode: 10001,
  },
  favorites: [
    "pizza", 
    42, 
    false, 
    null, 
    {
      item: "book",
      price: 10.99,
      in_stock: true,
    }
  ],
}
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
