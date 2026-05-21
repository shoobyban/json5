package json5

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Decode parses JSON5 data and stores the result in the value pointed to by v.
// Its signature mirrors encoding/json.Unmarshal: it accepts []byte input and
// decodes into a caller-supplied pointer.
//
// If v is nil or not a pointer, Decode returns an error.
//
// Supported target types:
//   - *any, *map[string]any, *[]any — populated directly from the parsed tree
//   - *string, *bool, *int, *float64 — populated when the JSON5 value is a matching scalar
//   - *struct — decoded via encoding/json round-trip (JSON5 → generic tree → JSON → struct)
//
// For struct targets, standard `json` struct tags are honoured via encoding/json.
func Decode(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("json5: Decode requires a non-nil pointer")
	}

	parsed, err := Unmarshal(string(data))
	if err != nil {
		return err
	}

	elem := rv.Elem()

	// Fast path: only empty interfaces can safely accept any parsed value.
	if elem.Kind() == reflect.Interface && elem.Type().NumMethod() == 0 {
		if parsed == nil {
			elem.Set(reflect.Zero(elem.Type()))
		} else {
			elem.Set(reflect.ValueOf(parsed))
		}
		return nil
	}

	// Fast path: direct type match for common untyped containers.
	switch target := v.(type) {
	case *map[string]any:
		m, ok := parsed.(map[string]any)
		if !ok {
			return fmt.Errorf("json5: cannot decode %T into *map[string]any", parsed)
		}
		*target = m
		return nil
	case *[]any:
		a, ok := parsed.([]any)
		if !ok {
			return fmt.Errorf("json5: cannot decode %T into *[]any", parsed)
		}
		*target = a
		return nil
	case *string:
		s, ok := parsed.(string)
		if !ok {
			return fmt.Errorf("json5: cannot decode %T into *string", parsed)
		}
		*target = s
		return nil
	case *bool:
		b, ok := parsed.(bool)
		if !ok {
			return fmt.Errorf("json5: cannot decode %T into *bool", parsed)
		}
		*target = b
		return nil
	case *int:
		n, ok := parsed.(int)
		if !ok {
			return fmt.Errorf("json5: cannot decode %T into *int", parsed)
		}
		*target = n
		return nil
	case *float64:
		f, ok := parsed.(float64)
		if !ok {
			// int → float64 promotion
			if n, ok2 := parsed.(int); ok2 {
				*target = float64(n)
				return nil
			}
			return fmt.Errorf("json5: cannot decode %T into *float64", parsed)
		}
		*target = f
		return nil
	}

	// Slow path for structs and other complex types: re-encode the generic
	// tree as standard JSON, then let encoding/json decode into the target.
	// This honours `json` struct tags automatically.
	jsonBytes, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("json5: intermediate marshal failed: %w", err)
	}
	if err := json.Unmarshal(jsonBytes, v); err != nil {
		return fmt.Errorf("json5: decode into %T failed: %w", v, err)
	}
	return nil
}

// Encode returns the JSON5 encoding of v as []byte.
// Its signature mirrors encoding/json.Marshal.
//
// Encode currently supports the same value types as [Marshal]:
// nil, bool, numeric types, string, []any, and map[string]any.
func Encode(v any) ([]byte, error) {
	s, err := Marshal(v)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// EncodeIndent returns the indented JSON5 encoding of v as []byte.
// Its signature mirrors encoding/json.MarshalIndent.
//
// The prefix string is prepended to each line of output (for encoding/json
// compatibility); indent is the string used for each level of indentation.
func EncodeIndent(v any, prefix, indent string) ([]byte, error) {
	s, err := MarshalIndent(v, indent)
	if err != nil {
		return nil, err
	}
	if prefix != "" {
		// Apply prefix to each line after the first, matching encoding/json behaviour.
		lines := splitLines(s)
		for i := 1; i < len(lines); i++ {
			lines[i] = prefix + lines[i]
		}
		s = joinLines(lines)
	}
	return []byte(s), nil
}

// Valid reports whether data is valid JSON5.
func Valid(data []byte) bool {
	_, err := Unmarshal(string(data))
	return err == nil
}

// splitLines splits s on newline boundaries, keeping the newlines at the end of each segment.
func splitLines(s string) []string {
	var lines []string
	for {
		i := 0
		for i < len(s) && s[i] != '\n' {
			i++
		}
		if i < len(s) {
			lines = append(lines, s[:i+1])
			s = s[i+1:]
		} else {
			lines = append(lines, s)
			break
		}
	}
	return lines
}

// joinLines concatenates the line segments back into a single string.
func joinLines(lines []string) string {
	n := 0
	for _, l := range lines {
		n += len(l)
	}
	buf := make([]byte, 0, n)
	for _, l := range lines {
		buf = append(buf, l...)
	}
	return string(buf)
}
