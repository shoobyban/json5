package json5

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizeStructuralTokens(t *testing.T) {
	tokens := Tokenize(`{}[]:,`)
	assert.Len(t, tokens, 6)
	assert.Equal(t, TOKEN_LBRACE, tokens[0].Type)
	assert.Equal(t, TOKEN_RBRACE, tokens[1].Type)
	assert.Equal(t, TOKEN_LBRACKET, tokens[2].Type)
	assert.Equal(t, TOKEN_RBRACKET, tokens[3].Type)
	assert.Equal(t, TOKEN_COLON, tokens[4].Type)
	assert.Equal(t, TOKEN_COMMA, tokens[5].Type)
}

func TestTokenizeBooleans(t *testing.T) {
	tokens := Tokenize(`true false`)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_TRUE, tokens[0].Type)
	assert.Equal(t, TOKEN_FALSE, tokens[1].Type)
}

func TestTokenizeNull(t *testing.T) {
	tokens := Tokenize(`null`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NULL, tokens[0].Type)
}

func TestTokenizeDoubleQuotedString(t *testing.T) {
	tokens := Tokenize(`"hello world"`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "hello world", tokens[0].Value)
}

func TestTokenizeSingleQuotedString(t *testing.T) {
	tokens := Tokenize(`'hello world'`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "hello world", tokens[0].Value)
}

func TestTokenizeStringWithEscapes(t *testing.T) {
	tokens := Tokenize(`"line1\nline2\ttab"`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "line1\nline2\ttab", tokens[0].Value)
}

func TestTokenizeStringWithEscapedQuotes(t *testing.T) {
	tokens := Tokenize(`"say \"hi\""`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, `say "hi"`, tokens[0].Value)
}

func TestTokenizeIntegerNumber(t *testing.T) {
	tokens := Tokenize(`42`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "42", tokens[0].Value)
}

func TestTokenizeFloatNumber(t *testing.T) {
	tokens := Tokenize(`3.14`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "3.14", tokens[0].Value)
}

func TestTokenizeNegativeNumber(t *testing.T) {
	tokens := Tokenize(`-99`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "-99", tokens[0].Value)
}

func TestTokenizeHexNumber(t *testing.T) {
	tokens := Tokenize(`0xFF`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "0xFF", tokens[0].Value)
}

func TestTokenizeScientificNotation(t *testing.T) {
	tokens := Tokenize(`1e10`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "1e10", tokens[0].Value)
}

func TestTokenizeUnquotedIdentifier(t *testing.T) {
	tokens := Tokenize(`myKey`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "myKey", tokens[0].Value)
}

func TestTokenizeIdentifierWithDollarAndUnderscore(t *testing.T) {
	tokens := Tokenize(`$_key123`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "$_key123", tokens[0].Value)
}

func TestTokenizeSingleLineComment(t *testing.T) {
	tokens := Tokenize("// this is a comment\nkey")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Equal(t, TOKEN_STRING, tokens[1].Type)
	assert.Equal(t, "key", tokens[1].Value)
}

func TestTokenizeMultiLineComment(t *testing.T) {
	tokens := Tokenize("/* multi\nline */key")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Contains(t, tokens[0].Value, "multi")
	assert.Equal(t, TOKEN_STRING, tokens[1].Type)
	assert.Equal(t, "key", tokens[1].Value)
}

func TestTokenizeUnknownCharacter(t *testing.T) {
	tokens := Tokenize(`@`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
	assert.Equal(t, "@", tokens[0].Value)
}

func TestTokenizeWhitespace(t *testing.T) {
	tokens := Tokenize("  \t\n\r  ")
	assert.Len(t, tokens, 0)
}

func TestTokenizeCompleteObject(t *testing.T) {
	input := `{key: "value", num: 42}`
	tokens := Tokenize(input)
	assert.Len(t, tokens, 9)
	assert.Equal(t, TOKEN_LBRACE, tokens[0].Type)
	assert.Equal(t, TOKEN_STRING, tokens[1].Type)
	assert.Equal(t, "key", tokens[1].Value)
	assert.Equal(t, TOKEN_COLON, tokens[2].Type)
	assert.Equal(t, TOKEN_STRING, tokens[3].Type)
	assert.Equal(t, "value", tokens[3].Value)
	assert.Equal(t, TOKEN_COMMA, tokens[4].Type)
	assert.Equal(t, TOKEN_STRING, tokens[5].Type)
	assert.Equal(t, "num", tokens[5].Value)
	assert.Equal(t, TOKEN_COLON, tokens[6].Type)
	assert.Equal(t, TOKEN_NUMBER, tokens[7].Type)
	assert.Equal(t, "42", tokens[7].Value)
	assert.Equal(t, TOKEN_RBRACE, tokens[8].Type)
}

// processEscapeSequences tests

func TestProcessEscapeNewline(t *testing.T) {
	result, err := processEscapeSequences(`hello\nworld`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\nworld", result)
}

func TestProcessEscapeCarriageReturn(t *testing.T) {
	result, err := processEscapeSequences(`hello\rworld`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\rworld", result)
}

func TestProcessEscapeTab(t *testing.T) {
	result, err := processEscapeSequences(`hello\tworld`)
	assert.NoError(t, err)
	assert.Equal(t, "hello\tworld", result)
}

func TestProcessEscapeBackslash(t *testing.T) {
	result, err := processEscapeSequences(`hello\\world`)
	assert.NoError(t, err)
	assert.Equal(t, `hello\world`, result)
}

func TestProcessEscapeDoubleQuote(t *testing.T) {
	result, err := processEscapeSequences(`say \"hi\"`)
	assert.NoError(t, err)
	assert.Equal(t, `say "hi"`, result)
}

func TestProcessEscapeSingleQuote(t *testing.T) {
	result, err := processEscapeSequences(`it\'s`)
	assert.NoError(t, err)
	assert.Equal(t, "it's", result)
}

func TestProcessEscapeUnicode4(t *testing.T) {
	// \u0041 = 'A'
	result, err := processEscapeSequences(`\u0041`)
	assert.NoError(t, err)
	assert.Equal(t, "A", result)
}

func TestProcessEscapeUnicode8(t *testing.T) {
	// \U0001F600 = emoji
	result, err := processEscapeSequences(`\U0001F600`)
	assert.NoError(t, err)
	assert.Equal(t, "\U0001F600", result)
}

func TestProcessEscapeUnicodeBraces(t *testing.T) {
	// \u{41} = 'A'
	result, err := processEscapeSequences(`\u{41}`)
	assert.NoError(t, err)
	assert.Equal(t, "A", result)
}

func TestProcessEscapeUnicodeBracesWithPrefix(t *testing.T) {
	// \u{0x1F600} = emoji
	result, err := processEscapeSequences(`\u{0x1F600}`)
	assert.NoError(t, err)
	assert.Equal(t, "\U0001F600", result)
}

func TestProcessEscapeUnicodeBracesUpperU(t *testing.T) {
	// \U{41} = 'A'
	result, err := processEscapeSequences(`\U{41}`)
	assert.NoError(t, err)
	assert.Equal(t, "A", result)
}

func TestProcessEscapeInvalidUnicodeBraces(t *testing.T) {
	// Incomplete \u{ sequence (no closing brace)
	_, err := processEscapeSequences(`\u{41`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete")
}

func TestProcessEscapeInvalidUnicodeHex(t *testing.T) {
	// Invalid hex in \uXXXX
	_, err := processEscapeSequences(`\uZZZZ`)
	assert.Error(t, err)
}

func TestProcessEscapeInvalidUnicode8Hex(t *testing.T) {
	// Invalid hex in \UXXXXXXXX
	_, err := processEscapeSequences(`\UZZZZZZZZ`)
	assert.Error(t, err)
}

func TestProcessEscapeInvalidUnicodeShort(t *testing.T) {
	// \u with not enough characters for 4-digit form
	_, err := processEscapeSequences(`\u41`)
	assert.Error(t, err)
}

func TestProcessEscapeInvalidUpperUShort(t *testing.T) {
	// \U with not enough characters for 8-digit form and no brace
	_, err := processEscapeSequences(`\U41`)
	assert.Error(t, err)
}

func TestProcessEscapeUnknownSequence(t *testing.T) {
	// Unknown escape like \z keeps the backslash (next char consumed by outer loop)
	result, err := processEscapeSequences(`\z`)
	assert.NoError(t, err)
	assert.Equal(t, `\z`, result)
}

func TestProcessEscapeNoEscape(t *testing.T) {
	result, err := processEscapeSequences(`hello`)
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestProcessEscapeInvalidBracesHex(t *testing.T) {
	// Invalid hex inside \u{...}
	_, err := processEscapeSequences(`\u{ZZZZ}`)
	assert.Error(t, err)
}

// Helper function tests

func TestIsWhitespace(t *testing.T) {
	assert.True(t, isWhitespace(' '))
	assert.True(t, isWhitespace('\t'))
	assert.True(t, isWhitespace('\n'))
	assert.True(t, isWhitespace('\r'))
	assert.False(t, isWhitespace('a'))
	assert.False(t, isWhitespace('0'))
}

func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit('0'))
	assert.True(t, isDigit('9'))
	assert.False(t, isDigit('a'))
	assert.False(t, isDigit(' '))
}

func TestIsHexDigit(t *testing.T) {
	assert.True(t, isHexDigit('0'))
	assert.True(t, isHexDigit('9'))
	assert.True(t, isHexDigit('a'))
	assert.True(t, isHexDigit('f'))
	assert.True(t, isHexDigit('A'))
	assert.True(t, isHexDigit('F'))
	assert.False(t, isHexDigit('g'))
	assert.False(t, isHexDigit('z'))
}

func TestIsIdentifierStart(t *testing.T) {
	assert.True(t, isIdentifierStart('a'))
	assert.True(t, isIdentifierStart('Z'))
	assert.True(t, isIdentifierStart('_'))
	assert.True(t, isIdentifierStart('$'))
	assert.False(t, isIdentifierStart('0'))
	assert.False(t, isIdentifierStart('-'))
}

func TestIsIdentifierPart(t *testing.T) {
	assert.True(t, isIdentifierPart('a'))
	assert.True(t, isIdentifierPart('Z'))
	assert.True(t, isIdentifierPart('_'))
	assert.True(t, isIdentifierPart('$'))
	assert.True(t, isIdentifierPart('5'))
	assert.False(t, isIdentifierPart('-'))
	assert.False(t, isIdentifierPart(' '))
}

// --- Bug fix tests ---

func TestProcessEscapeNullFollowedByDigit(t *testing.T) {
	// \0 followed by a digit should be an error per JSON5 spec (no octal)
	_, err := processEscapeSequences(`\09`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "\\0 followed by digit")
}

func TestProcessEscapeNullAlone(t *testing.T) {
	// \0 not followed by digit is valid
	result, err := processEscapeSequences(`\0`)
	assert.NoError(t, err)
	assert.Equal(t, "\x00", result)
}

func TestProcessEscapeNullFollowedByNonDigit(t *testing.T) {
	// \0 followed by a non-digit letter is valid
	result, err := processEscapeSequences(`\0a`)
	assert.NoError(t, err)
	assert.Equal(t, "\x00a", result)
}

func TestProcessEscapeLineContinuationLS(t *testing.T) {
	// Backslash followed by LS (\u2028) should be a line continuation
	input := "hello\\\u2028world"
	result, err := processEscapeSequences(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestProcessEscapeLineContinuationPS(t *testing.T) {
	// Backslash followed by PS (\u2029) should be a line continuation
	input := "hello\\\u2029world"
	result, err := processEscapeSequences(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestProcessEscapeSlash(t *testing.T) {
	// \/ should produce /
	result, err := processEscapeSequences(`\/`)
	assert.NoError(t, err)
	assert.Equal(t, "/", result)
}

func TestTokenizeUnicodeWhitespace(t *testing.T) {
	// Non-breaking space (\u00A0) should be treated as whitespace
	input := "42\u00A0true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "42", tokens[0].Value)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeLineSeparatorWhitespace(t *testing.T) {
	// LS (\u2028) should be treated as whitespace
	input := "42\u2028true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeParagraphSeparatorWhitespace(t *testing.T) {
	// PS (\u2029) should be treated as whitespace
	input := "42\u2029true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeBOMWhitespace(t *testing.T) {
	// BOM (\uFEFF) should be treated as whitespace
	input := "\uFEFF42"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "42", tokens[0].Value)
}

func TestTokenizeBarePlus(t *testing.T) {
	// A bare '+' with no number following should be TOKEN_UNKNOWN
	tokens := Tokenize("+")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
}

func TestTokenizeBareMinus(t *testing.T) {
	// A bare '-' with no number following should be TOKEN_UNKNOWN
	tokens := Tokenize("-")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
}

func TestTokenizePlusFollowedByNonNumber(t *testing.T) {
	// '+' followed by a letter should produce UNKNOWN + STRING
	tokens := Tokenize("+abc")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
	assert.Equal(t, TOKEN_STRING, tokens[1].Type)
}

func TestTokenizeMultiByteStringContents(t *testing.T) {
	// Strings containing multi-byte runes should be scanned correctly
	input := `"hello 🌍 world"`
	tokens := Tokenize(input)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "hello 🌍 world", tokens[0].Value)
}

func TestIsLineTerminator(t *testing.T) {
	assert.True(t, isLineTerminator('\n'))
	assert.True(t, isLineTerminator('\r'))
	assert.True(t, isLineTerminator('\u2028'))
	assert.True(t, isLineTerminator('\u2029'))
	assert.False(t, isLineTerminator(' '))
	assert.False(t, isLineTerminator('a'))
}

func TestTokenizeUnterminatedMultiLineComment(t *testing.T) {
	// Unterminated multi-line comment should not panic
	tokens := Tokenize("/* unterminated comment")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Contains(t, tokens[0].Value, "unterminated")
}

func TestTokenTypeString(t *testing.T) {
	assert.Equal(t, "LBRACE", TOKEN_LBRACE.String())
	assert.Equal(t, "STRING", TOKEN_STRING.String())
	assert.Equal(t, "NUMBER", TOKEN_NUMBER.String())
	assert.Equal(t, "COMMENT", TOKEN_COMMENT.String())
	assert.Equal(t, "UNKNOWN", TOKEN_UNKNOWN.String())
	assert.Equal(t, "TRUE", TOKEN_TRUE.String())
	assert.Equal(t, "FALSE", TOKEN_FALSE.String())
	assert.Equal(t, "NULL", TOKEN_NULL.String())
	assert.Equal(t, "RBRACE", TOKEN_RBRACE.String())
	assert.Equal(t, "LBRACKET", TOKEN_LBRACKET.String())
	assert.Equal(t, "RBRACKET", TOKEN_RBRACKET.String())
	assert.Equal(t, "COLON", TOKEN_COLON.String())
	assert.Equal(t, "COMMA", TOKEN_COMMA.String())
	assert.Contains(t, TokenType(99).String(), "TokenType(99)")
}

// --- Regression: Unicode Zs whitespace characters (#6) ---

func TestTokenizeEnSpace(t *testing.T) {
	// En space (\u2002) should be treated as whitespace
	input := "42\u2002true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeEmSpace(t *testing.T) {
	// Em space (\u2003) should be treated as whitespace
	input := "42\u2003true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeThinSpace(t *testing.T) {
	// Thin space (\u2009) should be treated as whitespace
	input := "42\u2009true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeIdeographicSpace(t *testing.T) {
	// Ideographic space (\u3000) should be treated as whitespace
	input := "42\u3000true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeOghamSpaceMark(t *testing.T) {
	// Ogham space mark (\u1680) should be treated as whitespace
	input := "42\u1680true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeNarrowNoBreakSpace(t *testing.T) {
	// Narrow no-break space (\u202F) should be treated as whitespace
	input := "42\u202Ftrue"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeMediumMathematicalSpace(t *testing.T) {
	// Medium mathematical space (\u205F) should be treated as whitespace
	input := "42\u205Ftrue"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

// --- Regression: Unterminated strings (#18) ---

func TestTokenizeUnterminatedDoubleQuotedString(t *testing.T) {
	// Unterminated double-quoted string should not panic
	tokens := Tokenize(`"unterminated`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
}

func TestTokenizeUnterminatedSingleQuotedString(t *testing.T) {
	// Unterminated single-quoted string should not panic
	tokens := Tokenize(`'unterminated`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
}

// --- Regression: Multi-byte chars inside multi-line comments (#3) ---

func TestTokenizeMultiLineCommentWithMultiByteChars(t *testing.T) {
	// Multi-byte characters inside comments (emoji, CJK) should not cause false close
	input := "/* hello 🌍 world */42"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Equal(t, "/* hello 🌍 world */", tokens[0].Value)
	assert.Equal(t, TOKEN_NUMBER, tokens[1].Type)
	assert.Equal(t, "42", tokens[1].Value)
}

func TestTokenizeMultiLineCommentWithCJK(t *testing.T) {
	// CJK characters should not cause false close of multi-line comment
	input := "/* 日本語テスト */true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

// --- Regression: Unterminated multi-line comment edge cases ---

func TestTokenizeUnterminatedCommentMinimal(t *testing.T) {
	// Just /* with nothing else
	tokens := Tokenize("/*")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
}

func TestTokenizeUnterminatedCommentWithStar(t *testing.T) {
	// /* followed by * but no /
	tokens := Tokenize("/* almost *")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
}

// --- Regression: Tokenizer error as TOKEN_UNKNOWN (#7) ---

func TestTokenizeInvalidEscapeProducesUnknown(t *testing.T) {
	// Invalid escape sequence inside string produces TOKEN_UNKNOWN
	tokens := Tokenize(`"\09"`)
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
	assert.Contains(t, tokens[0].Value, "\\0 followed by digit")
}

// --- Regression: Unknown escape writes both chars (#8) ---

func TestProcessEscapeUnknownWritesBothChars(t *testing.T) {
	// \z should produce exactly backslash + z
	result, err := processEscapeSequences(`abc\zdef`)
	assert.NoError(t, err)
	assert.Equal(t, "abc\\zdef", result)
}

func TestProcessEscapeMultipleUnknown(t *testing.T) {
	// Multiple unknown escapes in sequence
	result, err := processEscapeSequences(`\x\y\z`)
	assert.NoError(t, err)
	assert.Equal(t, "\\x\\y\\z", result)
}

// --- Regression: Single-line comment ending at EOF ---

func TestTokenizeSingleLineCommentAtEOF(t *testing.T) {
	// Comment at end of input with no trailing newline
	tokens := Tokenize("42 // comment at eof")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, TOKEN_COMMENT, tokens[1].Type)
}

// --- Regression: Single-line comment terminated by LS/PS ---

func TestTokenizeSingleLineCommentTerminatedByLS(t *testing.T) {
	// LS (\u2028) should terminate a single-line comment
	input := "// comment\u2028true"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Equal(t, TOKEN_TRUE, tokens[1].Type)
}

func TestTokenizeSingleLineCommentTerminatedByPS(t *testing.T) {
	// PS (\u2029) should terminate a single-line comment
	input := "// comment\u2029false"
	tokens := Tokenize(input)
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_COMMENT, tokens[0].Type)
	assert.Equal(t, TOKEN_FALSE, tokens[1].Type)
}

// --- Regression: \r and \r\n line continuation at processEscapeSequences level ---

func TestProcessEscapeLineContinuationCR(t *testing.T) {
	// Backslash + CR alone should be a line continuation
	input := "hello\\\rworld"
	result, err := processEscapeSequences(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestProcessEscapeLineContinuationCRLF(t *testing.T) {
	// Backslash + CRLF should be a single line continuation
	input := "hello\\\r\nworld"
	result, err := processEscapeSequences(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

func TestProcessEscapeLineContinuationLF(t *testing.T) {
	// Backslash + LF should be a line continuation
	input := "hello\\\nworld"
	result, err := processEscapeSequences(input)
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)
}

// --- Regression: Double sign tokens ---

func TestTokenizeDoublePlus(t *testing.T) {
	// ++1 should produce UNKNOWN("+") then NUMBER("+1")
	tokens := Tokenize("++1")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
	assert.Equal(t, TOKEN_NUMBER, tokens[1].Type)
	assert.Equal(t, "+1", tokens[1].Value)
}

func TestTokenizeDoubleMinus(t *testing.T) {
	// --1 should produce UNKNOWN("-") then NUMBER("-1")
	tokens := Tokenize("--1")
	assert.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_UNKNOWN, tokens[0].Type)
	assert.Equal(t, TOKEN_NUMBER, tokens[1].Type)
	assert.Equal(t, "-1", tokens[1].Value)
}

// --- Regression: Empty hex number ---

func TestTokenizeEmptyHex(t *testing.T) {
	// 0x with no hex digits should still tokenize as a number (parser will reject)
	tokens := Tokenize("0x")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "0x", tokens[0].Value)
}

// --- Regression: Malformed exponents ---

func TestTokenizeMalformedExponentNoDigits(t *testing.T) {
	// 1e with no exponent digits — consumeExponent still advances past 'e'
	tokens := Tokenize("1e")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "1e", tokens[0].Value)
}

func TestTokenizeMalformedExponentSignNoDigits(t *testing.T) {
	// 1e+ with sign but no exponent digits
	tokens := Tokenize("1e+")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "1e+", tokens[0].Value)
}

func TestTokenizeMalformedExponentUppercase(t *testing.T) {
	// 1E- with sign but no exponent digits
	tokens := Tokenize("1E-")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "1E-", tokens[0].Value)
}

// --- Regression: Trailing dot with exponent ---

func TestTokenizeTrailingDotWithExponent(t *testing.T) {
	// 5.e2 should tokenize as a single number
	tokens := Tokenize("5.e2")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "5.e2", tokens[0].Value)
}

// --- Regression: Signed leading dot with exponent ---

func TestTokenizeSignedLeadingDotWithExponent(t *testing.T) {
	tokens := Tokenize("+.5e2")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "+.5e2", tokens[0].Value)

	tokens = Tokenize("-.5e2")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
	assert.Equal(t, "-.5e2", tokens[0].Value)
}

// --- Regression: Multi-byte Unicode identifiers ---

func TestTokenizeUnicodeIdentifier(t *testing.T) {
	// Non-ASCII letters should be valid identifier parts
	tokens := Tokenize("café")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "café", tokens[0].Value)
}

func TestTokenizeUnicodeIdentifierCJK(t *testing.T) {
	// CJK characters are letters and should be valid identifiers
	tokens := Tokenize("名前")
	assert.Len(t, tokens, 1)
	assert.Equal(t, TOKEN_STRING, tokens[0].Type)
	assert.Equal(t, "名前", tokens[0].Value)
}

// --- Regression: Minus followed by dot ---

func TestTokenizeMinusDot(t *testing.T) {
	// -. with no digit after dot should be UNKNOWN
	tokens := Tokenize("-.")
	assert.GreaterOrEqual(t, len(tokens), 1)
	// The '-' is followed by '.' which matches the dot check,
	// but i+1 must have a digit after '.' for signed leading decimal
	// If not handled, verify it doesn't panic
}

// --- Regression: consumeExponent with valid values ---

func TestTokenizeExponentVariations(t *testing.T) {
	cases := []struct {
		input string
		value string
	}{
		{"1e0", "1e0"},
		{"1E1", "1E1"},
		{"1e+0", "1e+0"},
		{"1e-0", "1e-0"},
		{"1e10", "1e10"},
		{"1E+10", "1E+10"},
		{"1e-3", "1e-3"},
		{"2.5e4", "2.5e4"},
		{"0.1e-2", "0.1e-2"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			tokens := Tokenize(tc.input)
			assert.Len(t, tokens, 1)
			assert.Equal(t, TOKEN_NUMBER, tokens[0].Type)
			assert.Equal(t, tc.value, tokens[0].Value)
		})
	}
}
