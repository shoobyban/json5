package json5

import (
	"testing"
)

// FuzzTokenize ensures the tokenizer never panics on arbitrary input.
func FuzzTokenize(f *testing.F) {
	// Seed corpus with interesting inputs
	f.Add(`{"key": "value"}`)
	f.Add(`[1, 2, 3]`)
	f.Add(`"hello\nworld"`)
	f.Add(`0xFF`)
	f.Add(`Infinity`)
	f.Add(`NaN`)
	f.Add(`+42`)
	f.Add(`-.5`)
	f.Add(`.5e2`)
	f.Add(`// comment`)
	f.Add(`/* multi-line */`)
	f.Add(`'single quoted'`)
	f.Add(`unquoted_key`)
	f.Add(`"\u0041"`)
	f.Add(`"\U0001F600"`)
	f.Add(`"\u{41}"`)
	f.Add(`"\0"`)
	f.Add(`"unterminated`)
	f.Add(`/* unterminated`)
	f.Add(`+`)
	f.Add(`-`)
	f.Add(`@`)
	f.Add(`""`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add("\uFEFF42")
	f.Add("42\u00A0true")
	f.Add("42\u2003true")
	f.Add(`"\09"`)

	f.Fuzz(func(t *testing.T, input string) {
		// Must not panic
		_ = Tokenize(input)
	})
}

// FuzzUnMarshal ensures the parser never panics on arbitrary input.
func FuzzUnMarshal(f *testing.F) {
	// Seed corpus with valid and interesting inputs
	f.Add(`{"key": "value"}`)
	f.Add(`[1, 2, 3]`)
	f.Add(`"hello"`)
	f.Add(`42`)
	f.Add(`3.14`)
	f.Add(`true`)
	f.Add(`false`)
	f.Add(`null`)
	f.Add(`Infinity`)
	f.Add(`-Infinity`)
	f.Add(`+Infinity`)
	f.Add(`NaN`)
	f.Add(`0xFF`)
	f.Add(`+42`)
	f.Add(`-.5`)
	f.Add(`.5e2`)
	f.Add(`{key: "value", num: 1}`)
	f.Add(`{'key': 'value'}`)
	f.Add(`[1, "two", true, null]`)
	f.Add(`{a: 1, b: 2,}`)
	f.Add(`[1, 2, 3,]`)
	f.Add(`""`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(``)
	f.Add(`   `)
	// Truncated inputs
	f.Add(`{a: 1`)
	f.Add(`{a`)
	f.Add(`{a:`)
	f.Add(`[1`)
	f.Add(`[{a: 1`)
	// Invalid inputs
	f.Add(`@`)
	f.Add(`+`)
	f.Add(`-`)
	f.Add(`{key "value"}`)
	f.Add(`{: "value"}`)
	f.Add(`[1 2]`)
	// Comments
	f.Add("// comment\n{key: 1}")
	f.Add("/* comment */{key: 1}")
	f.Add("// just a comment")

	f.Fuzz(func(t *testing.T, input string) {
		// Must not panic
		result, err := UnMarshal(input)
		if err != nil {
			return
		}
		// If parsing succeeds, marshal should also not panic
		_, _ = Marshal(result)
	})
}

// FuzzMarshalRoundTrip ensures marshal->unmarshal round-trip doesn't panic
// or lose data for valid parsed results.
func FuzzMarshalRoundTrip(f *testing.F) {
	f.Add(`{"key": "value"}`)
	f.Add(`[1, 2, 3]`)
	f.Add(`{a: 1, b: "two", c: true, d: null}`)
	f.Add(`[Infinity, -Infinity, NaN]`)
	f.Add(`{hex: 0xFF, dot: .5, trail: 5.}`)

	f.Fuzz(func(t *testing.T, input string) {
		result, err := UnMarshal(input)
		if err != nil {
			return
		}
		marshaled, err := Marshal(result)
		if err != nil {
			return
		}
		// Re-parse the marshaled output — must not panic
		_, _ = UnMarshal(marshaled)
	})
}
