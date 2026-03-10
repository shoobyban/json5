package json5

import (
	"math"
	"strings"
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

	result, err := Unmarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 object: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}

	expected := map[string]any{
		"address":   map[string]any{"city": "New York", "zipcode": 10001},
		"age":       42,
		"children":  nil,
		"favorites": []any{"pizza", 42, false, nil, map[string]any{"in_stock": true, "item": "book", "price": 10.99}},
		"married":   true,
		"name":      "John Doe", "hexadecimal": 912559, "lineBreaks": "line 1\nline 2\n\t\tline 3",
		"unicode": "Hello beauty! -😀-"}
	assert.Equal(t, expected, result)
}

func TestParseJSON5Array(t *testing.T) {
	input := `["a", 'b', 1]`

	result, err := Unmarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 array: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func TestParseJSON5String(t *testing.T) {
	input := `'a'`

	result, err := Unmarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 string: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func TestParseJSON5Number(t *testing.T) {
	input := `42`

	result, err := Unmarshal(input)
	if err != nil {
		t.Fatalf("Error parsing JSON5 number: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}
}

func TestParseEmptyInput(t *testing.T) {
	result, err := Unmarshal("")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseWhitespaceOnly(t *testing.T) {
	result, err := Unmarshal("   \t\n  ")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseBoolTrue(t *testing.T) {
	result, err := Unmarshal("true")
	assert.NoError(t, err)
	assert.Equal(t, true, result)
}

func TestParseBoolFalse(t *testing.T) {
	result, err := Unmarshal("false")
	assert.NoError(t, err)
	assert.Equal(t, false, result)
}

func TestParseNull(t *testing.T) {
	result, err := Unmarshal("null")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseFloat(t *testing.T) {
	result, err := Unmarshal("3.14")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, result)
}

func TestParseNegativeNumber(t *testing.T) {
	result, err := Unmarshal("-7")
	assert.NoError(t, err)
	assert.Equal(t, -7, result)
}

func TestParseNegativeFloat(t *testing.T) {
	result, err := Unmarshal("-2.5")
	assert.NoError(t, err)
	assert.Equal(t, -2.5, result)
}

func TestParseHexNumber(t *testing.T) {
	result, err := Unmarshal("0xFF")
	assert.NoError(t, err)
	assert.Equal(t, 255, result)
}

func TestParseObjectWithComments(t *testing.T) {
	input := `{
		// single-line comment
		key: "value",
	}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", obj["key"])
}

func TestParseObjectWithMultiLineComment(t *testing.T) {
	input := `{
		/* multi
		   line comment */
		key: "value",
	}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", obj["key"])
}

func TestUnmarshalAlias(t *testing.T) {
	result, err := Unmarshal(`{"key": "value"}`)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", obj["key"])
}

func TestParseObjectWithoutComments(t *testing.T) {
	input := `{key: "value", num: 1}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", obj["key"])
	assert.Equal(t, 1, obj["num"])
}

func TestParseEmptyObject(t *testing.T) {
	result, err := Unmarshal("{}")
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Len(t, obj, 0)
}

func TestParseEmptyArray(t *testing.T) {
	result, err := Unmarshal("[]")
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Len(t, arr, 0)
}

func TestParseNestedArray(t *testing.T) {
	input := `[[1, 2], [3, 4]]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
	inner, ok := arr[0].([]any)
	assert.True(t, ok)
	assert.Equal(t, []any{1, 2}, inner)
}

func TestParseObjectMissingColon(t *testing.T) {
	input := `{key "value"}`
	_, err := Unmarshal(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected ':'")
}

func TestParseObjectBadToken(t *testing.T) {
	// A non-string token where a key is expected
	input := `{: "value"}`
	_, err := Unmarshal(input)
	assert.Error(t, err)
}

func TestParseObjectMissingCommaOrBrace(t *testing.T) {
	input := `{key: "value" other: "val"}`
	_, err := Unmarshal(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected ',' or '}'")
}

func TestParseArrayMissingCommaOrBracket(t *testing.T) {
	input := `[1 2]`
	_, err := Unmarshal(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected ',' or ']'")
}

func TestParseUnexpectedToken(t *testing.T) {
	// A lone unexpected character at top level
	input := `@`
	_, err := Unmarshal(input)
	assert.Error(t, err)
}

func TestParseValueUnexpectedToken(t *testing.T) {
	// Inside an array, an unexpected comma-as-value scenario
	input := `{key: }`
	_, err := Unmarshal(input)
	assert.Error(t, err)
}

func TestParseObjectWithSingleQuoteKeys(t *testing.T) {
	input := `{'key': 'value'}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", obj["key"])
}

func TestParseArrayWithMixedTypes(t *testing.T) {
	input := `[1, "two", true, null, 3.14, {key: "val"}, [5]]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Len(t, arr, 7)
	assert.Equal(t, 1, arr[0])
	assert.Equal(t, "two", arr[1])
	assert.Equal(t, true, arr[2])
	assert.Nil(t, arr[3])
	assert.Equal(t, 3.14, arr[4])
}

func TestParseScientificNotation(t *testing.T) {
	result, err := Unmarshal("1e3")
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, result)
}

func TestParseObjectTrailingComma(t *testing.T) {
	input := `{a: 1, b: 2,}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, 1, obj["a"])
	assert.Equal(t, 2, obj["b"])
}

func TestParseArrayTrailingComma(t *testing.T) {
	input := `[1, 2, 3,]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Equal(t, []any{1, 2, 3}, arr)
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
		Unmarshal(input)
	}
}

// --- JSON5 Number Features ---

func TestParseInfinity(t *testing.T) {
	result, err := Unmarshal(`Infinity`)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(1), result)
}

func TestParseNegativeInfinity(t *testing.T) {
	result, err := Unmarshal(`-Infinity`)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(-1), result)
}

func TestParsePositiveInfinity(t *testing.T) {
	result, err := Unmarshal(`+Infinity`)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(1), result)
}

func TestParseNaN(t *testing.T) {
	result, err := Unmarshal(`NaN`)
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(result.(float64)))
}

func TestParseLeadingDecimalPoint(t *testing.T) {
	result, err := Unmarshal(`.5`)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, result)
}

func TestParseTrailingDecimalPoint(t *testing.T) {
	result, err := Unmarshal(`5.`)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, result)
}

func TestParsePositiveSign(t *testing.T) {
	result, err := Unmarshal(`+42`)
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestParsePositiveFloat(t *testing.T) {
	result, err := Unmarshal(`+3.14`)
	assert.NoError(t, err)
	assert.Equal(t, 3.14, result)
}

func TestParseExponentWithSign(t *testing.T) {
	result, err := Unmarshal(`1e+2`)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, result)

	result, err = Unmarshal(`1e-2`)
	assert.NoError(t, err)
	assert.Equal(t, 0.01, result)

	result, err = Unmarshal(`1E+3`)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, result)
}

func TestParseNegativeHex(t *testing.T) {
	result, err := Unmarshal(`-0xFF`)
	assert.NoError(t, err)
	assert.Equal(t, -255, result)
}

func TestParsePositiveHex(t *testing.T) {
	result, err := Unmarshal(`+0xFF`)
	assert.NoError(t, err)
	assert.Equal(t, 255, result)
}

func TestParseSignedLeadingDot(t *testing.T) {
	result, err := Unmarshal(`+.5`)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, result)

	result, err = Unmarshal(`-.5`)
	assert.NoError(t, err)
	assert.Equal(t, -0.5, result)
}

func TestParseLeadingDotWithExponent(t *testing.T) {
	result, err := Unmarshal(`.5e2`)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, result)
}

// --- JSON5 String Features ---

func TestParseEscapeBackspace(t *testing.T) {
	result, err := Unmarshal(`"hello\bworld"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\bworld", result)
}

func TestParseEscapeFormFeed(t *testing.T) {
	result, err := Unmarshal(`"hello\fworld"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\fworld", result)
}

func TestParseEscapeVerticalTab(t *testing.T) {
	result, err := Unmarshal(`"hello\vworld"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\vworld", result)
}

func TestParseEscapeNullChar(t *testing.T) {
	result, err := Unmarshal(`"hello\0world"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\x00world", result)
}

func TestParseLineContinuation(t *testing.T) {
	// Backslash followed by newline should produce empty string (line continuation)
	input := "\"hello\\\nworld\""
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestParseLineContinuationCRLF(t *testing.T) {
	// Backslash followed by CR LF should produce empty string (line continuation)
	input := "\"hello\\\r\nworld\""
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

// --- JSON5 Special Values in Objects/Arrays ---

func TestParseObjectWithSpecialNumbers(t *testing.T) {
	input := `{
		inf: Infinity,
		negInf: -Infinity,
		nan: NaN,
		hex: 0xCAFE,
		negHex: -0xBEEF,
		leadingDot: .5,
		trailingDot: 5.,
		pos: +42,
		exp: 1e+5,
	}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)

	obj := result.(map[string]any)
	assert.Equal(t, math.Inf(1), obj["inf"])
	assert.Equal(t, math.Inf(-1), obj["negInf"])
	assert.True(t, math.IsNaN(obj["nan"].(float64)))
	assert.Equal(t, 0xCAFE, obj["hex"])
	assert.Equal(t, -0xBEEF, obj["negHex"])
	assert.Equal(t, 0.5, obj["leadingDot"])
	assert.Equal(t, 5.0, obj["trailingDot"])
	assert.Equal(t, 42, obj["pos"])
	assert.Equal(t, 100000.0, obj["exp"])
}

func TestParseArrayWithSpecialNumbers(t *testing.T) {
	input := `[Infinity, -Infinity, NaN, .5, 5., +3, 1e+2]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)

	arr := result.([]any)
	assert.Equal(t, math.Inf(1), arr[0])
	assert.Equal(t, math.Inf(-1), arr[1])
	assert.True(t, math.IsNaN(arr[2].(float64)))
	assert.Equal(t, 0.5, arr[3])
	assert.Equal(t, 5.0, arr[4])
	assert.Equal(t, 3, arr[5])
	assert.Equal(t, 100.0, arr[6])
}

// --- Edge case / negative tests ---

func TestParseBarePlus(t *testing.T) {
	_, err := Unmarshal("+")
	assert.Error(t, err)
}

func TestParseBareMinus(t *testing.T) {
	_, err := Unmarshal("-")
	assert.Error(t, err)
}

func TestParseNullEscapeFollowedByDigit(t *testing.T) {
	// \0 followed by a digit should produce an error
	_, err := Unmarshal(`"\09"`)
	assert.Error(t, err)
}

func TestParseLineContinuationLS(t *testing.T) {
	// Backslash followed by LS (\u2028) should be a line continuation
	input := "\"hello\\\u2028world\""
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestParseLineContinuationPS(t *testing.T) {
	// Backslash followed by PS (\u2029) should be a line continuation
	input := "\"hello\\\u2029world\""
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestParseUnicodeWhitespace(t *testing.T) {
	// Non-breaking space as separator between tokens
	input := "{a:\u00A01}"
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, 1, obj["a"])
}

// --- Truncated input tests ---

func TestParseTruncatedObject(t *testing.T) {
	_, err := Unmarshal(`{a: 1`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}

func TestParseTruncatedObjectAfterKey(t *testing.T) {
	_, err := Unmarshal(`{a`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}

func TestParseTruncatedArray(t *testing.T) {
	_, err := Unmarshal(`[1`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}

// --- Regression: Additional truncated input variants (#1) ---

func TestParseTruncatedObjectAfterColon(t *testing.T) {
	// Object truncated right after colon — parseValue sees end of input
	_, err := Unmarshal(`{a:`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}

func TestParseTruncatedNestedObject(t *testing.T) {
	// Nested structure truncated mid-parse
	_, err := Unmarshal(`[{a: 1`)
	assert.Error(t, err)
}

func TestParseTruncatedArrayInObject(t *testing.T) {
	_, err := Unmarshal(`{a: [1, 2`)
	assert.Error(t, err)
}

func TestParseTruncatedObjectInArray(t *testing.T) {
	_, err := Unmarshal(`[{a: 1`)
	assert.Error(t, err)
}

func TestParseTruncatedDeeplyNested(t *testing.T) {
	_, err := Unmarshal(`{a: {b: {c:`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of input")
}

// --- Regression: Comments in arrays (#4) ---

func TestParseArrayWithSingleLineComment(t *testing.T) {
	input := `[
		// comment
		1,
		2,
	]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Equal(t, []any{1, 2}, arr)
}

func TestParseArrayWithMultiLineComment(t *testing.T) {
	input := `[
		/* multi-line
		   comment */
		1,
		2,
	]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr, ok := result.([]any)
	assert.True(t, ok)
	assert.Equal(t, []any{1, 2}, arr)
}

func TestParseCommentBetweenKeyAndColon(t *testing.T) {
	input := `{key /* comment */ : "value"}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, "value", obj["key"])
}

func TestParseCommentBetweenColonAndValue(t *testing.T) {
	input := `{key: /* comment */ "value"}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, "value", obj["key"])
}

func TestParseCommentBetweenArrayElements(t *testing.T) {
	input := `[1, /* inline */ 2, 3]`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	arr := result.([]any)
	assert.Equal(t, []any{1, 2, 3}, arr)
}

func TestParseCommentOnlyInput(t *testing.T) {
	result, err := Unmarshal("// just a comment")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseMultipleCommentsOnlyInput(t *testing.T) {
	result, err := Unmarshal("// comment 1\n/* comment 2 */")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// --- Regression: Deeply nested structures (#18) ---

func TestParseDeeplyNestedObjects(t *testing.T) {
	// Build a deeply nested object: {a: {a: {a: ... "leaf" ...}}}
	depth := 100
	var sb strings.Builder
	for i := 0; i < depth; i++ {
		sb.WriteString(`{a: `)
	}
	sb.WriteString(`"leaf"`)
	for i := 0; i < depth; i++ {
		sb.WriteString(`}`)
	}
	result, err := Unmarshal(sb.String())
	assert.NoError(t, err)
	// Walk down to the leaf
	current := result
	for i := 0; i < depth; i++ {
		obj, ok := current.(map[string]any)
		assert.True(t, ok, "level %d should be a map", i)
		current = obj["a"]
	}
	assert.Equal(t, "leaf", current)
}

func TestParseDeeplyNestedArrays(t *testing.T) {
	// Build deeply nested arrays: [[[... 42 ...]]]
	depth := 100
	var sb strings.Builder
	for i := 0; i < depth; i++ {
		sb.WriteString(`[`)
	}
	sb.WriteString(`42`)
	for i := 0; i < depth; i++ {
		sb.WriteString(`]`)
	}
	result, err := Unmarshal(sb.String())
	assert.NoError(t, err)
	current := result
	for i := 0; i < depth; i++ {
		arr, ok := current.([]any)
		assert.True(t, ok, "level %d should be an array", i)
		assert.Len(t, arr, 1)
		current = arr[0]
	}
	assert.Equal(t, 42, current)
}

// --- Regression: Unterminated strings through parser (#18) ---

func TestParseUnterminatedString(t *testing.T) {
	// Parser should handle unterminated string without panic
	// The tokenizer produces a TOKEN_STRING even if unterminated
	result, err := Unmarshal(`"unterminated`)
	// Should not panic — either returns the partial string or an error
	if err == nil {
		assert.Equal(t, "unterminated", result)
	}
}

// --- Regression: Unicode Zs whitespace in parser (#6) ---

func TestParseEnSpaceWhitespace(t *testing.T) {
	input := "{a:\u20021}"
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, 1, obj["a"])
}

func TestParseIdeographicSpaceWhitespace(t *testing.T) {
	input := "{a:\u30001}"
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, 1, obj["a"])
}

// --- Regression: Invalid hex numbers (#15) ---

func TestParseInvalidHexNumber(t *testing.T) {
	_, err := Unmarshal(`0xZZZZ`)
	assert.Error(t, err)
}

func TestParseInvalidNumber(t *testing.T) {
	// 123abc is tokenized as number(123) + string(abc).
	// The parser only consumes the first top-level token, so it returns 123.
	result, err := Unmarshal(`123abc`)
	assert.NoError(t, err)
	assert.Equal(t, 123, result)
}

// --- Regression: Round-trip with comments (#4) ---

func TestRoundTripObjectWithComments(t *testing.T) {
	input := `{
		// name field
		name: "Alice",
		/* age */
		age: 30,
	}`
	parsed, err := Unmarshal(input)
	assert.NoError(t, err)
	marshaled, err := Marshal(parsed)
	assert.NoError(t, err)
	reparsed, err := Unmarshal(marshaled)
	assert.NoError(t, err)
	assert.Equal(t, parsed, reparsed)
}

// --- Regression: Double sign produces error ---

func TestParseDoublePlus(t *testing.T) {
	_, err := Unmarshal("++1")
	assert.Error(t, err)
}

func TestParseDoubleMinus(t *testing.T) {
	_, err := Unmarshal("--1")
	assert.Error(t, err)
}

// --- Regression: Empty hex number ---

func TestParseEmptyHex(t *testing.T) {
	_, err := Unmarshal("0x")
	assert.Error(t, err)
}

// --- Regression: Malformed exponents ---

func TestParseMalformedExponentNoDigits(t *testing.T) {
	// 1e with no digits after exponent
	_, err := Unmarshal("1e")
	assert.Error(t, err)
}

func TestParseMalformedExponentSignOnly(t *testing.T) {
	// 1e+ with no digits after the sign
	_, err := Unmarshal("1e+")
	assert.Error(t, err)
}

// --- Regression: Trailing dot with exponent ---

func TestParseTrailingDotWithExponent(t *testing.T) {
	result, err := Unmarshal("5.e2")
	assert.NoError(t, err)
	assert.Equal(t, 500.0, result)
}

// --- Regression: Signed leading dot with exponent ---

func TestParseSignedLeadingDotWithExponent(t *testing.T) {
	result, err := Unmarshal("+.5e2")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, result)

	result, err = Unmarshal("-.5e2")
	assert.NoError(t, err)
	assert.Equal(t, -50.0, result)
}

// --- Regression: Multi-byte Unicode key in object ---

func TestParseUnicodeKeyObject(t *testing.T) {
	input := `{café: 42}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, 42, obj["café"])
}

func TestParseCJKKeyObject(t *testing.T) {
	input := `{名前: "太郎"}`
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	obj := result.(map[string]any)
	assert.Equal(t, "太郎", obj["名前"])
}

// --- Regression: \r line continuation through parser ---

func TestParseLineContinuationCR(t *testing.T) {
	// Backslash + lone CR should be a line continuation
	input := "\"hello\\\rworld\""
	result, err := Unmarshal(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

// --- Regression: Slash escape round-trip ---

func TestParseSlashEscape(t *testing.T) {
	result, err := Unmarshal(`"hello\/world"`)
	assert.NoError(t, err)
	assert.Equal(t, "hello/world", result)
}
