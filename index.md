# JSON5 for Go

A lightweight JSON5 parser and marshaller for Go with full TinyGo support.

## What is JSON5?

JSON5 is a superset of JSON that makes configuration files more human-friendly. It adds features like:

- **Comments** - Single-line (`//`) and multi-line (`/* */`)
- **Unquoted keys** - Write `name: "value"` instead of `"name": "value"`
- **Trailing commas** - No more syntax errors from that last comma
- **Hexadecimal numbers** - Use `0xdecaf` for readability
- **Multi-line strings** - Break long strings across lines
- **Single quotes** - Use `'value'` or `"value"` interchangeably

## Installation

```bash
go get github.com/shoobyban/json5
```

## Quick Start

### Parsing JSON5

```go
package main

import (
    "fmt"
    "github.com/shoobyban/json5"
)

func main() {
    input := `{
        // This is a comment
        name: "Alice",
        age: 30,
        active: true,
    }`
    
    result, err := json5.Unmarshal(input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("%+v\n", result)
}
```

### Creating JSON5

```go
package main

import (
    "fmt"
    "github.com/shoobyban/json5"
)

func main() {
    data := map[string]any{
        "name": "Alice",
        "age": 30,
        "active": true,
    }
    
    output, err := json5.MarshalIndent(data, "  ")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

## encoding/json-Compatible API

For users migrating from `encoding/json`, a compatibility layer provides
familiar signatures that accept `[]byte` and decode into typed Go values
(including structs with `json` tags).

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

### Encode / EncodeIndent (like json.Marshal / json.MarshalIndent)

```go
data := map[string]any{"name": "Alice", "active": true}

b, err := json5.Encode(data)
// b == []byte(`{active: true, name: "Alice"}`)

b, err = json5.EncodeIndent(data, "", "  ")
```

### Valid (like json.Valid)

```go
json5.Valid([]byte(`{key: "value",}`))  // true
json5.Valid([]byte(`{key: }`))           // false
```

## API Reference

### encoding/json-Compatible Functions

#### `Decode(data []byte, v any) error`

Parses JSON5 data and stores the result in the value pointed to by v.
Supports `*any`, `*map[string]any`, `*[]any`, scalar pointers, and
structs (via `encoding/json` round-trip honouring `json` struct tags).

#### `Encode(v any) ([]byte, error)`

Returns the JSON5 encoding of v as `[]byte`.

#### `EncodeIndent(v any, prefix, indent string) ([]byte, error)`

Returns the indented JSON5 encoding of v as `[]byte`.

#### `Valid(data []byte) bool`

Reports whether data is valid JSON5.

### Original Functions

#### `Unmarshal(input string) (any, error)`

Parses a JSON5 string and returns the result as a Go value.

**Returns:**

- `map[string]any` for objects
- `[]any` for arrays
- `string` for strings
- `int` or `float64` for numbers
- `bool` for booleans
- `nil` for null

**Example:**

```go
result, err := json5.Unmarshal(`{name: "Bob", age: 25}`)
if err != nil {
    // Handle error
}

person := result.(map[string]any)
name := person["name"].(string)
age := person["age"].(int)
```

#### `UnMarshal(input string) (any, error)` — Deprecated

> **Deprecated:** Use `Unmarshal` (lowercase 'm') or `Decode` instead.

Alias for `Unmarshal`. Will continue to work but may be removed in a
future major version.

### Marshalling Functions

#### `Marshal(value any) (string, error)`

Converts a Go value to JSON5 format without indentation.

**Example:**
```go
data := map[string]any{
    "city": "Tokyo",
    "population": 14000000,
}

output, err := json5.Marshal(data)
// Output: {city: "Tokyo", population: 14000000}
```

#### `MarshalIndent(value any, indent string) (string, error)`

Converts a Go value to JSON5 format with pretty-printing.

**Parameters:**
- `value` - The data to marshal
- `indent` - Indentation string (e.g., `"  "` or `"\t"`)

**Example:**
```go
data := map[string]any{
    "users": []any{
        map[string]any{"name": "Alice"},
        map[string]any{"name": "Bob"},
    },
}

output, err := json5.MarshalIndent(data, "  ")
/*
Output:
{
  users: [
    {
      name: "Alice",
    }, 
    {
      name: "Bob",
    }
  ],
}
*/
```

### Tokenizer

For advanced use cases, you can use the tokenizer directly.

#### `Tokenize(input string) []Token`

Breaks a JSON5 string into tokens.

**Example:**
```go
tokens := json5.Tokenize(`{name: "Alice", age: 30}`)
for _, token := range tokens {
    fmt.Printf("Type: %d, Value: %s\n", token.Type, token.Value)
}
```

## Supported Features

### ✅ Comments

```json5
{
    // Single-line comment
    name: "Alice",
    
    /* Multi-line
       comment */
    age: 30
}
```

### ✅ Unquoted Object Keys

```json5
{
    firstName: "Alice",
    lastName: "Smith",
    age: 30
}
```

### ✅ Trailing Commas

```json5
{
    name: "Alice",
    age: 30,  // This comma is fine
}
```

### ✅ Single and Double Quotes

```json5
{
    single: 'value',
    double: "value"
}
```

### ✅ Hexadecimal Numbers

```json5
{
    color: 0xFF5733,
    flag: 0xdecaf
}
```

### ✅ Escape Sequences

```json5
{
    newline: "line 1\nline 2",
    tab: "column1\tcolumn2",
    backslash: "path\\to\\file"
}
```

### ✅ Unicode Escapes

```json5
{
    emoji: "Hello \u{1F600}",  // 😀
    symbol: "\U{0x2764}"        // ❤
}
```

### ✅ Multi-line Strings

```json5
{
    description: "This is a long string
        that spans multiple lines"
}
```

## Type Conversion

When parsing JSON5, values are converted to Go types:

| JSON5 Type | Go Type |
|------------|---------|
| Object | `map[string]any` |
| Array | `[]any` |
| String | `string` |
| Number | `int` or `float64` |
| Boolean | `bool` |
| Null | `nil` |

### Working with Parsed Data

```go
result, _ := json5.Unmarshal(`{
    name: "Alice",
    age: 30,
    scores: [95, 87, 92]
}`)

// Type assertions
person := result.(map[string]any)
name := person["name"].(string)
age := person["age"].(int)
scores := person["scores"].([]any)

fmt.Printf("%s is %d years old\n", name, age)
```

## TinyGo Support

This library is designed to work with TinyGo, which doesn't support reflection. All parsing and marshalling uses explicit type handling instead of reflection-based conversion.

```go
// Works in both Go and TinyGo
result, err := json5.Unmarshal(configString)
config := result.(map[string]any)
```

## Error Handling

The library returns descriptive errors for invalid JSON5:

```go
result, err := json5.Unmarshal(`{invalid json5}`)
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
    return
}
```

## Common Use Cases

### Configuration Files

```go
// Load and parse a JSON5 config file
configData := `{
    // Server configuration
    server: {
        host: "localhost",
        port: 8080,
    },
    
    // Database settings
    database: {
        host: "localhost",
        port: 5432,
        name: "myapp",
    }
}`

config, _ := json5.Unmarshal(configData)
```

### API Responses

```go
// Create a JSON5 response
response := map[string]any{
    "success": true,
    "data": map[string]any{
        "id": 123,
        "name": "Alice",
    },
}

output, _ := json5.MarshalIndent(response, "  ")
```

### Data Transformation

```go
// Parse JSON5, modify, and output
input := `{name: "Alice", age: 30}`
data, _ := json5.Unmarshal(input)

person := data.(map[string]any)
person["age"] = person["age"].(int) + 1

output, _ := json5.Marshal(person)
```

## Limitations

- Performance is not optimized for extremely large files
- Input is string-based (not streaming)
- Integers are parsed as `int`, decimals as `float64`

For most configuration files and small to medium data processing, performance is sufficient.

## Contributing

Pull requests are welcome! The library is intentionally simple and focused on TinyGo compatibility.

## License

MIT License - see LICENSE file for details.

## Links

- [JSON5 Specification](https://json5.org/)
- [TinyGo Project](https://tinygo.org/)
- [GitHub Repository](https://github.com/shoobyban/json5)
