package json5

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	titanous "github.com/titanous/json5"
	furukawa "github.com/yosuke-furukawa/json5/encoding/json5"
)

// generateStdJSON creates a standard JSON object with n key-value pairs.
func generateStdJSON(n int) string {
	var b strings.Builder
	b.WriteString(`{`)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(`,`)
		}
		fmt.Fprintf(&b, `"key%d": "value%d"`, i, i)
	}
	// add nested object and array
	fmt.Fprintf(&b, `,"nested": {"a": 1, "b": 2, "c": [1, 2, 3]}`)
	fmt.Fprintf(&b, `,"flag": true, "nothing": null, "pi": 3.14159`)
	b.WriteString(`}`)
	return b.String()
}

// generateJSON5 creates a JSON5 object with n key-value pairs using JSON5 features
// that are commonly supported: unquoted keys, trailing commas, single-line comments.
func generateJSON5(n int) string {
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString("  // this is a JSON5 object\n")
	for i := 0; i < n; i++ {
		// unquoted keys, trailing comma
		fmt.Fprintf(&b, "  key%d: \"value%d\",\n", i, i)
	}
	// nested object with unquoted keys
	b.WriteString("  nested: {a: 1, b: 2, c: [1, 2, 3,],},\n")
	b.WriteString("  flag: true,\n")
	b.WriteString("  nothing: null,\n")
	b.WriteString("  pi: 3.14159,\n")
	b.WriteString("}")
	return b.String()
}

// --- Standard JSON benchmarks (all 3 parsers + encoding/json baseline) ---

func BenchmarkStdJSON_Small_Builtin(b *testing.B) {
	input := []byte(generateStdJSON(5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		json.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Small_Shoobyban(b *testing.B) {
	input := generateStdJSON(5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkStdJSON_Small_Titanous(b *testing.B) {
	input := []byte(generateStdJSON(5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Small_Furukawa(b *testing.B) {
	input := []byte(generateStdJSON(5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Medium_Builtin(b *testing.B) {
	input := []byte(generateStdJSON(50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		json.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Medium_Shoobyban(b *testing.B) {
	input := generateStdJSON(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkStdJSON_Medium_Titanous(b *testing.B) {
	input := []byte(generateStdJSON(50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Medium_Furukawa(b *testing.B) {
	input := []byte(generateStdJSON(50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Large_Builtin(b *testing.B) {
	input := []byte(generateStdJSON(500))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		json.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Large_Shoobyban(b *testing.B) {
	input := generateStdJSON(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkStdJSON_Large_Titanous(b *testing.B) {
	input := []byte(generateStdJSON(500))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkStdJSON_Large_Furukawa(b *testing.B) {
	input := []byte(generateStdJSON(500))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}

// --- JSON5 benchmarks (3 JSON5 parsers only, no encoding/json) ---

func BenchmarkJSON5_Small_Shoobyban(b *testing.B) {
	input := generateJSON5(5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkJSON5_Small_Titanous(b *testing.B) {
	input := []byte(generateJSON5(5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkJSON5_Small_Furukawa(b *testing.B) {
	input := []byte(generateJSON5(5))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}

func BenchmarkJSON5_Medium_Shoobyban(b *testing.B) {
	input := generateJSON5(50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkJSON5_Medium_Titanous(b *testing.B) {
	input := []byte(generateJSON5(50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkJSON5_Medium_Furukawa(b *testing.B) {
	input := []byte(generateJSON5(50))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}

func BenchmarkJSON5_Large_Shoobyban(b *testing.B) {
	input := generateJSON5(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(input)
	}
}

func BenchmarkJSON5_Large_Titanous(b *testing.B) {
	input := []byte(generateJSON5(500))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		titanous.Unmarshal(input, &v)
	}
}

func BenchmarkJSON5_Large_Furukawa(b *testing.B) {
	input := []byte(generateJSON5(500))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v interface{}
		furukawa.Unmarshal(input, &v)
	}
}
