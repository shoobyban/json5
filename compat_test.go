package json5

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Decode tests
// ---------------------------------------------------------------------------

func TestDecodeIntoAny(t *testing.T) {
	var v any
	err := Decode([]byte(`{key: "value", num: 42}`), &v)
	assert.NoError(t, err)
	m, ok := v.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", m["key"])
	assert.Equal(t, 42, m["num"])
}

func TestDecodeIntoAnyNull(t *testing.T) {
	var v any = "sentinel"
	err := Decode([]byte(`null`), &v)
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestDecodeIntoAnyEmptyInput(t *testing.T) {
	var v any = "sentinel"
	err := Decode([]byte(``), &v)
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestDecodeIntoMap(t *testing.T) {
	var m map[string]any
	err := Decode([]byte(`{a: 1, b: "two"}`), &m)
	assert.NoError(t, err)
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, "two", m["b"])
}

func TestDecodeIntoSlice(t *testing.T) {
	var s []any
	err := Decode([]byte(`[1, "two", true]`), &s)
	assert.NoError(t, err)
	assert.Equal(t, []any{1, "two", true}, s)
}

func TestDecodeIntoString(t *testing.T) {
	var s string
	err := Decode([]byte(`"hello"`), &s)
	assert.NoError(t, err)
	assert.Equal(t, "hello", s)
}

func TestDecodeIntoBool(t *testing.T) {
	var b bool
	err := Decode([]byte(`true`), &b)
	assert.NoError(t, err)
	assert.True(t, b)
}

func TestDecodeIntoInt(t *testing.T) {
	var n int
	err := Decode([]byte(`42`), &n)
	assert.NoError(t, err)
	assert.Equal(t, 42, n)
}

func TestDecodeIntoFloat64(t *testing.T) {
	var f float64
	err := Decode([]byte(`3.14`), &f)
	assert.NoError(t, err)
	assert.Equal(t, 3.14, f)
}

func TestDecodeIntPromotesToFloat64(t *testing.T) {
	var f float64
	err := Decode([]byte(`42`), &f)
	assert.NoError(t, err)
	assert.Equal(t, 42.0, f)
}

func TestDecodeIntoStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var p Person
	err := Decode([]byte(`{name: "Alice", age: 30}`), &p)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", p.Name)
	assert.Equal(t, 30, p.Age)
}

func TestDecodeIntoStructNested(t *testing.T) {
	type Address struct {
		City string `json:"city"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}
	var p Person
	err := Decode([]byte(`{name: "Bob", address: {city: "NYC"}}`), &p)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", p.Name)
	assert.Equal(t, "NYC", p.Address.City)
}

func TestDecodeIntoStructWithOmittedFields(t *testing.T) {
	type Config struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Name string `json:"name"`
	}
	var c Config
	// JSON5 only has host — Port and Name should be zero-valued.
	err := Decode([]byte(`{host: "localhost"}`), &c)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", c.Host)
	assert.Equal(t, 0, c.Port)
	assert.Equal(t, "", c.Name)
}

func TestDecodeIntoStructWithArray(t *testing.T) {
	type Data struct {
		Tags []string `json:"tags"`
	}
	var d Data
	err := Decode([]byte(`{tags: ["go", "json5"]}`), &d)
	assert.NoError(t, err)
	assert.Equal(t, []string{"go", "json5"}, d.Tags)
}

func TestDecodeNilPointer(t *testing.T) {
	err := Decode([]byte(`{}`), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-nil pointer")
}

func TestDecodeNonPointer(t *testing.T) {
	var m map[string]any
	err := Decode([]byte(`{}`), m) // not &m
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-nil pointer")
}

func TestDecodeTypeMismatch(t *testing.T) {
	var s string
	err := Decode([]byte(`42`), &s)
	assert.Error(t, err)
}

func TestDecodeMapTypeMismatch(t *testing.T) {
	var m map[string]any
	err := Decode([]byte(`[1, 2]`), &m)
	assert.Error(t, err)
}

func TestDecodeSliceTypeMismatch(t *testing.T) {
	var s []any
	err := Decode([]byte(`{a: 1}`), &s)
	assert.Error(t, err)
}

func TestDecodeBoolTypeMismatch(t *testing.T) {
	var b bool
	err := Decode([]byte(`"hello"`), &b)
	assert.Error(t, err)
}

func TestDecodeIntTypeMismatch(t *testing.T) {
	var n int
	err := Decode([]byte(`"hello"`), &n)
	assert.Error(t, err)
}

func TestDecodeFloat64TypeMismatchString(t *testing.T) {
	var f float64
	err := Decode([]byte(`"hello"`), &f)
	assert.Error(t, err)
}

func TestDecodeInvalidJSON5(t *testing.T) {
	var v any
	err := Decode([]byte(`{key: }`), &v)
	assert.Error(t, err)
}

func TestDecodeRejectsTrailingTokens(t *testing.T) {
	var v any
	err := Decode([]byte(`42 true`), &v)
	assert.Error(t, err)
}

func TestDecodeJSON5Features(t *testing.T) {
	// Decode should handle all JSON5 features: comments, unquoted keys,
	// trailing commas, hex numbers, single-quoted strings.
	input := []byte(`{
		// comment
		name: 'Alice',
		hex: 0xFF,
		active: true,
	}`)
	var m map[string]any
	err := Decode(input, &m)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", m["name"])
	assert.Equal(t, 255, m["hex"])
	assert.Equal(t, true, m["active"])
}

func TestDecodeIntoFloat64Infinity(t *testing.T) {
	var f float64
	err := Decode([]byte(`Infinity`), &f)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(1), f)
}

func TestDecodeIntoFloat64NegInfinity(t *testing.T) {
	var f float64
	err := Decode([]byte(`-Infinity`), &f)
	assert.NoError(t, err)
	assert.Equal(t, math.Inf(-1), f)
}

func TestDecodeIntoFloat64NaN(t *testing.T) {
	var f float64
	err := Decode([]byte(`NaN`), &f)
	assert.NoError(t, err)
	assert.True(t, math.IsNaN(f))
}

// ---------------------------------------------------------------------------
// Encode tests
// ---------------------------------------------------------------------------

func TestEncodeBasic(t *testing.T) {
	result, err := Encode(map[string]any{"key": "value"})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{key: "value"}`), result)
}

func TestEncodeNil(t *testing.T) {
	result, err := Encode(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), result)
}

func TestEncodeReturnsBytes(t *testing.T) {
	result, err := Encode(42)
	assert.NoError(t, err)
	assert.IsType(t, []byte{}, result)
	assert.Equal(t, []byte("42"), result)
}

func TestEncodeString(t *testing.T) {
	result, err := Encode("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"hello"`), result)
}

func TestEncodeBool(t *testing.T) {
	result, err := Encode(true)
	assert.NoError(t, err)
	assert.Equal(t, []byte("true"), result)
}

func TestEncodeArray(t *testing.T) {
	result, err := Encode([]any{1, "two", false})
	assert.NoError(t, err)
	assert.Equal(t, []byte(`[1, "two", false]`), result)
}

func TestEncodeUnsupported(t *testing.T) {
	_, err := Encode(make(chan int))
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// EncodeIndent tests
// ---------------------------------------------------------------------------

func TestEncodeIndentBasic(t *testing.T) {
	input := map[string]any{"a": 1}
	result, err := EncodeIndent(input, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, []byte("{\n  a: 1,\n}"), result)
}

func TestEncodeIndentWithPrefix(t *testing.T) {
	input := map[string]any{"a": 1}
	result, err := EncodeIndent(input, ">> ", "  ")
	assert.NoError(t, err)
	expected := "{\n>> " + "  a: 1,\n>> }"
	assert.Equal(t, []byte(expected), result)
}

func TestEncodeIndentEmpty(t *testing.T) {
	result, err := EncodeIndent(map[string]any{}, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, []byte("{}"), result)
}

func TestEncodeIndentNoPrefix(t *testing.T) {
	input := map[string]any{"x": []any{1, 2}}
	result, err := EncodeIndent(input, "", "\t")
	assert.NoError(t, err)
	assert.Contains(t, string(result), "\n")
	assert.Contains(t, string(result), "\t")
}

func TestEncodeIndentUnsupported(t *testing.T) {
	_, err := EncodeIndent(make(chan int), "", "  ")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// Valid tests
// ---------------------------------------------------------------------------

func TestValidTrue(t *testing.T) {
	cases := []string{
		`{key: "value"}`,
		`[1, 2, 3]`,
		`"hello"`,
		`42`,
		`true`,
		`null`,
		`Infinity`,
		`NaN`,
		`{a: 1, b: 2,}`,   // trailing comma
		`{/* comment */}`,  // comment
		`{'key': 'value'}`, // single quotes
		`{hex: 0xFF}`,      // hex number
	}
	for _, c := range cases {
		assert.True(t, Valid([]byte(c)), "expected valid: %s", c)
	}
}

func TestValidFalseWithTrailingTokens(t *testing.T) {
	assert.False(t, Valid([]byte(`42 true`)))
}

func TestValidFalse(t *testing.T) {
	cases := []string{
		`{key: }`,
		`[1 2]`,
		`{key "value"}`,
	}
	for _, c := range cases {
		assert.False(t, Valid([]byte(c)), "expected invalid: %s", c)
	}
}

// ---------------------------------------------------------------------------
// Round-trip: Encode then Decode
// ---------------------------------------------------------------------------

func TestRoundTripEncodeDecodeMap(t *testing.T) {
	original := map[string]any{"name": "Alice", "age": 30}
	encoded, err := Encode(original)
	assert.NoError(t, err)

	var decoded map[string]any
	err = Decode(encoded, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", decoded["name"])
	assert.Equal(t, 30, decoded["age"])
}

func TestRoundTripEncodeDecodeSlice(t *testing.T) {
	original := []any{"a", 1, true, nil}
	encoded, err := Encode(original)
	assert.NoError(t, err)

	var decoded []any
	err = Decode(encoded, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "a", decoded[0])
	assert.Equal(t, 1, decoded[1])
	assert.Equal(t, true, decoded[2])
	assert.Nil(t, decoded[3])
}

func TestRoundTripEncodeDecodeStruct(t *testing.T) {
	type Item struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}
	original := map[string]any{"name": "Widget", "price": 10}
	encoded, err := Encode(original)
	assert.NoError(t, err)

	var item Item
	err = Decode(encoded, &item)
	assert.NoError(t, err)
	assert.Equal(t, "Widget", item.Name)
	assert.Equal(t, 10, item.Price)
}

func TestRoundTripEncodeDecodeScalar(t *testing.T) {
	encoded, err := Encode("hello world")
	assert.NoError(t, err)

	var s string
	err = Decode(encoded, &s)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", s)
}

// ---------------------------------------------------------------------------
// Backward compatibility: existing API unchanged
// ---------------------------------------------------------------------------

func TestExistingMarshalUnchanged(t *testing.T) {
	s, err := Marshal(map[string]any{"k": 1})
	assert.NoError(t, err)
	assert.IsType(t, "", s) // returns string, not []byte
	assert.Equal(t, "{k: 1}", s)
}

func TestExistingUnmarshalUnchanged(t *testing.T) {
	v, err := Unmarshal(`{k: 1}`)
	assert.NoError(t, err)
	assert.IsType(t, map[string]any{}, v) // returns any
	m := v.(map[string]any)
	assert.Equal(t, 1, m["k"])
}

func TestExistingMarshalIndentUnchanged(t *testing.T) {
	s, err := MarshalIndent(map[string]any{"k": 1}, "  ")
	assert.NoError(t, err)
	assert.IsType(t, "", s)
	assert.Contains(t, s, "\n")
}

func TestDeprecatedUnMarshalStillWorks(t *testing.T) {
	v, err := UnMarshal(`{k: 1}`)
	assert.NoError(t, err)
	m := v.(map[string]any)
	assert.Equal(t, 1, m["k"])
}
