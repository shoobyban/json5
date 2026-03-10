package json5

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalNil(t *testing.T) {
	result, err := Marshal(nil)
	assert.NoError(t, err)
	assert.Equal(t, "null", result)
}

func TestMarshalBool(t *testing.T) {
	result, err := Marshal(true)
	assert.NoError(t, err)
	assert.Equal(t, "true", result)

	result, err = Marshal(false)
	assert.NoError(t, err)
	assert.Equal(t, "false", result)
}

func TestMarshalNumber(t *testing.T) {
	result, err := Marshal(42)
	assert.NoError(t, err)
	assert.Equal(t, "42", result)

	result, err = Marshal(3.14)
	assert.NoError(t, err)
	assert.Equal(t, "3.14", result)
}

func TestMarshalNumberTypes(t *testing.T) {
	// Cover all numeric type branches in marshalValue
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(1), "1"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"float32", float32(1.5), "1.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarshalString(t *testing.T) {
	result, err := Marshal("Hello\nWorld")
	assert.NoError(t, err)
	assert.Equal(t, "\"Hello\\nWorld\"", result)
}

func TestMarshalStringEscaping(t *testing.T) {
	// Covers all escape branches in marshalString
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"backslash", `back\slash`, `"back\\slash"`},
		{"double quote", `say "hi"`, `"say \"hi\""`},
		{"carriage return", "hello\rworld", `"hello\rworld"`},
		{"tab", "hello\tworld", `"hello\tworld"`},
		{"newline", "hello\nworld", `"hello\nworld"`},
		{"combined", "a\tb\nc\\d\"e\r", `"a\tb\nc\\d\"e\r"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarshalUnsupportedType(t *testing.T) {
	// A channel is not a supported marshal type
	ch := make(chan int)
	_, err := Marshal(ch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestMarshalArray(t *testing.T) {
	input := []any{"a", 1, true, nil}
	expected := `["a", 1, true, null]`
	result, err := Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalArrayEmpty(t *testing.T) {
	result, err := Marshal([]any{})
	assert.NoError(t, err)
	assert.Equal(t, "[]", result)
}

func TestMarshalObjectEmpty(t *testing.T) {
	result, err := Marshal(map[string]any{})
	assert.NoError(t, err)
	assert.Equal(t, "{}", result)
}

func TestMarshalArrayNestedError(t *testing.T) {
	// An unsupported type inside an array should propagate the error
	input := []any{"ok", make(chan int)}
	_, err := Marshal(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestMarshalObjectNestedError(t *testing.T) {
	// An unsupported type as a value in an object should propagate the error
	input := map[string]any{"bad": make(chan int)}
	_, err := Marshal(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestMarshalIndent(t *testing.T) {
	input := map[string]any{
		"name": "Alice",
		"pets": []any{"cat", "dog"},
	}
	expected := "{\n  name: \"Alice\",\n  pets: [\n    \"cat\",\n    \"dog\",\n  ],\n}"
	result, err := MarshalIndent(input, "  ")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalIndentNested(t *testing.T) {
	input := map[string]any{
		"outer": map[string]any{
			"inner": "value",
		},
	}
	expected := "{\n\touter: {\n\t\tinner: \"value\",\n\t},\n}"
	result, err := MarshalIndent(input, "\t")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshalKeyIdentifiers(t *testing.T) {
	// Keys that are simple identifiers should not be quoted
	// Keys that are not simple identifiers must be quoted
	tests := []struct {
		name     string
		input    map[string]any
		contains string
	}{
		{"dollar sign key", map[string]any{"$key": 1}, "$key: 1"},
		{"underscore key", map[string]any{"_key": 1}, "_key: 1"},
		{"alphanumeric key", map[string]any{"key1": 1}, "key1: 1"},
		{"empty key", map[string]any{"": 1}, `"": 1`},
		{"key with space", map[string]any{"a b": 1}, `"a b": 1`},
		{"key starts with digit", map[string]any{"1key": 1}, `"1key": 1`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			assert.NoError(t, err)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func TestMarshalObject(t *testing.T) {
	input := map[string]any{
		"name": "John Doe",
		"age":  42,
		"address": map[string]any{
			"city":    "New York",
			"zipcode": 10001,
		},
	}
	expected := `{address: {city: "New York", zipcode: 10001}, age: 42, name: "John Doe"}`
	result, err := Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// --- Marshal JSON5 Special Values ---

func TestMarshalInfinity(t *testing.T) {
	result, err := Marshal(math.Inf(1))
	assert.NoError(t, err)
	assert.Equal(t, "Infinity", result)
}

func TestMarshalNegativeInfinity(t *testing.T) {
	result, err := Marshal(math.Inf(-1))
	assert.NoError(t, err)
	assert.Equal(t, "-Infinity", result)
}

func TestMarshalNaN(t *testing.T) {
	result, err := Marshal(math.NaN())
	assert.NoError(t, err)
	assert.Equal(t, "NaN", result)
}

func TestMarshalFloat32Infinity(t *testing.T) {
	result, err := Marshal(float32(math.Inf(1)))
	assert.NoError(t, err)
	assert.Equal(t, "Infinity", result)
}

func TestMarshalFloat32NaN(t *testing.T) {
	result, err := Marshal(float32(math.NaN()))
	assert.NoError(t, err)
	assert.Equal(t, "NaN", result)
}

// --- Round-trip: marshal then unmarshal ---

func TestRoundTripInfinity(t *testing.T) {
	marshaled, err := Marshal(math.Inf(1))
	assert.NoError(t, err)

	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(1), result)
}

func TestRoundTripNegativeInfinity(t *testing.T) {
	marshaled, err := Marshal(math.Inf(-1))
	assert.NoError(t, err)

	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(-1), result)
}

func TestRoundTripNaN(t *testing.T) {
	marshaled, err := Marshal(math.NaN())
	assert.NoError(t, err)

	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(result.(float64)))
}

func TestMarshalComplexObject(t *testing.T) {
	input := map[string]any{
		"simpleKey": "value",
		"complex key": map[string]any{
			"nestedKey": "nestedValue",
		},
	}
	expected := `{"complex key": {nestedKey: "nestedValue"}, simpleKey: "value"}`
	result, err := Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

// --- marshalString control character escaping ---

func TestMarshalStringBackspace(t *testing.T) {
	result, err := Marshal("hello\bworld")
	assert.NoError(t, err)
	assert.Equal(t, `"hello\bworld"`, result)
}

func TestMarshalStringFormFeed(t *testing.T) {
	result, err := Marshal("hello\fworld")
	assert.NoError(t, err)
	assert.Equal(t, `"hello\fworld"`, result)
}

func TestMarshalStringVerticalTab(t *testing.T) {
	result, err := Marshal("hello\vworld")
	assert.NoError(t, err)
	assert.Equal(t, `"hello\vworld"`, result)
}

func TestMarshalStringNullByte(t *testing.T) {
	result, err := Marshal("hello\x00world")
	assert.NoError(t, err)
	assert.Equal(t, `"hello\0world"`, result)
}

// --- Round-trip tests for control characters ---

func TestRoundTripBackspace(t *testing.T) {
	original := "a\bb"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRoundTripFormFeed(t *testing.T) {
	original := "a\fb"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRoundTripVerticalTab(t *testing.T) {
	original := "a\vb"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestRoundTripNullByte(t *testing.T) {
	original := "a\x00b"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

// --- Regression: Compact vs Indented Marshal (#9, #10, #11) ---

func TestMarshalCompactVsIndentObject(t *testing.T) {
	input := map[string]any{"a": 1, "b": 2}

	compact, err := Marshal(input)
	assert.NoError(t, err)
	assert.NotContains(t, compact, "\n", "compact output should have no newlines")

	indented, err := MarshalIndent(input, "  ")
	assert.NoError(t, err)
	assert.Contains(t, indented, "\n", "indented output should have newlines")
	assert.Contains(t, indented, "  ", "indented output should have indent")
}

func TestMarshalCompactVsIndentArray(t *testing.T) {
	input := []any{1, 2, 3}

	compact, err := Marshal(input)
	assert.NoError(t, err)
	assert.NotContains(t, compact, "\n", "compact output should have no newlines")
	assert.Equal(t, "[1, 2, 3]", compact)

	indented, err := MarshalIndent(input, "\t")
	assert.NoError(t, err)
	assert.Contains(t, indented, "\n")
	assert.Contains(t, indented, "\t")
}

func TestMarshalIndentArrayTrailingCommas(t *testing.T) {
	// Indented arrays should have trailing commas on each element
	input := []any{"a", "b"}
	result, err := MarshalIndent(input, "  ")
	assert.NoError(t, err)
	assert.Equal(t, "[\n  \"a\",\n  \"b\",\n]", result)
}

func TestMarshalIndentObjectTrailingCommas(t *testing.T) {
	// Indented objects should have trailing commas on each entry
	input := map[string]any{"z": 1}
	result, err := MarshalIndent(input, "  ")
	assert.NoError(t, err)
	assert.Equal(t, "{\n  z: 1,\n}", result)
}

func TestMarshalCompactNoTrailingComma(t *testing.T) {
	// Compact arrays and objects should NOT have trailing commas
	arrResult, err := Marshal([]any{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, "[1, 2]", arrResult)

	objResult, err := Marshal(map[string]any{"a": 1})
	assert.NoError(t, err)
	assert.Equal(t, "{a: 1}", objResult)
}

// --- Regression: Round-trip full complex object ---

func TestRoundTripComplexObject(t *testing.T) {
	input := map[string]any{
		"name":   "Alice",
		"age":    30,
		"active": true,
		"score":  3.14,
		"data":   nil,
		"tags":   []any{"go", "json5"},
		"nested": map[string]any{
			"inner": "value",
		},
	}
	marshaled, err := Marshal(input)
	assert.NoError(t, err)
	parsed, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, input, parsed)
}

func TestRoundTripIndentedComplexObject(t *testing.T) {
	input := map[string]any{
		"items": []any{1, "two", true, nil},
		"meta":  map[string]any{"version": 2},
	}
	marshaled, err := MarshalIndent(input, "  ")
	assert.NoError(t, err)
	parsed, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, input, parsed)
}

// --- Regression: Marshal float32 negative infinity ---

func TestMarshalFloat32NegativeInfinity(t *testing.T) {
	result, err := Marshal(float32(math.Inf(-1)))
	assert.NoError(t, err)
	assert.Equal(t, "-Infinity", result)
}

// --- Regression: Round-trip string with slash ---

func TestRoundTripStringWithSlash(t *testing.T) {
	original := "hello/world"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

// --- Regression: Marshal string with all control chars combined ---

func TestMarshalStringAllControlChars(t *testing.T) {
	// String containing all escapable control characters
	original := "\b\f\v\x00\t\n\r"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	assert.Equal(t, `"\b\f\v\0\t\n\r"`, marshaled)

	// Round-trip
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

// --- Regression: Marshal/unmarshal multi-byte strings ---

func TestRoundTripMultiByteString(t *testing.T) {
	original := "hello 🌍 café 日本語"
	marshaled, err := Marshal(original)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

// --- Regression: Marshal/unmarshal special numbers in containers ---

func TestRoundTripSpecialNumbersInArray(t *testing.T) {
	input := []any{math.Inf(1), math.Inf(-1), math.NaN()}
	marshaled, err := Marshal(input)
	assert.NoError(t, err)
	result, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	arr := result.([]any)
	assert.Equal(t, math.Inf(1), arr[0])
	assert.Equal(t, math.Inf(-1), arr[1])
	assert.True(t, math.IsNaN(arr[2].(float64)))
}
