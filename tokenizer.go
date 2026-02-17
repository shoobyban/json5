package json5

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType represents different types of JSON5 tokens
type TokenType int

const (
	TOKEN_LBRACE   TokenType = iota // {
	TOKEN_RBRACE                    // }
	TOKEN_LBRACKET                  // [
	TOKEN_RBRACKET                  // ]
	TOKEN_COLON                     // :
	TOKEN_COMMA                     // ,
	TOKEN_STRING                    // string (quoted or unquoted)
	TOKEN_NUMBER                    // number (including hex)
	TOKEN_TRUE                      // true
	TOKEN_FALSE                     // false
	TOKEN_NULL                      // null
	TOKEN_COMMENT                   // comment
	TOKEN_UNKNOWN                   // unknown
)

// Token represents a JSON5 token with its type and value
type Token struct {
	Type  TokenType
	Value string
}

// String returns the name of the token type.
func (t TokenType) String() string {
	switch t {
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_COLON:
		return "COLON"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_TRUE:
		return "TRUE"
	case TOKEN_FALSE:
		return "FALSE"
	case TOKEN_NULL:
		return "NULL"
	case TOKEN_COMMENT:
		return "COMMENT"
	case TOKEN_UNKNOWN:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("TokenType(%d)", int(t))
	}
}

// isWhitespace checks if a character is a JSON5 whitespace character.
// JSON5 extends whitespace to include tab, vertical tab, form feed,
// space, non-breaking space, BOM, and any Unicode space separator,
// plus the line terminators LF, CR, LS (\u2028), and PS (\u2029).
func isWhitespace(ch rune) bool {
	switch ch {
	case '\t', '\n', '\r',
		'\v',     // vertical tab
		'\f',     // form feed
		'\uFEFF', // BOM / zero-width no-break space
		'\u2028', // line separator
		'\u2029': // paragraph separator
		return true
	}
	// Unicode category Zs (space separators) includes space, NBSP,
	// en/em space, thin space, and other Unicode spaces.
	return unicode.Is(unicode.Zs, ch)
}

// isLineTerminator checks if a rune is a JSON5 line terminator.
func isLineTerminator(ch rune) bool {
	return ch == '\n' || ch == '\r' || ch == '\u2028' || ch == '\u2029'
}

// isDigit checks if a character is a digit
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isHexDigit checks if a character is a valid hexadecimal digit
func isHexDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

// isIdentifierStart checks if a character can be the start of an unquoted key
func isIdentifierStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || ch == '$'
}

// isIdentifierPart checks if a character can be part of an unquoted key
func isIdentifierPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '$'
}

// processEscapeSequences converts escape sequences such as \n, \t, \b, \f, \v, \0,
// \uXXXX, \UXXXXXXXX, \u{0x1FA}, and \U{0x1FA} into their actual representations.
// Also handles line continuations (backslash followed by a line terminator including
// LF, CR, CRLF, LS \u2028, and PS \u2029).
func processEscapeSequences(input string) (string, error) {
	var result strings.Builder
	runes := []rune(input)
	length := len(runes)

	for i := 0; i < length; i++ {
		ch := runes[i]

		if ch == '\\' && i+1 < length {
			nextCh := runes[i+1]
			switch nextCh {
			case 'n':
				result.WriteByte('\n')
				i++
			case 'r':
				result.WriteByte('\r')
				i++
			case 't':
				result.WriteByte('\t')
				i++
			case '\\':
				result.WriteByte('\\')
				i++
			case '"':
				result.WriteByte('"')
				i++
			case '\'':
				result.WriteByte('\'')
				i++
			case '/':
				result.WriteByte('/')
				i++
			case 'b':
				result.WriteByte('\b')
				i++
			case 'f':
				result.WriteByte('\f')
				i++
			case 'v':
				result.WriteByte('\v')
				i++
			case '0':
				// \0 is only valid when NOT followed by another digit (to avoid octal ambiguity)
				if i+2 < length && isDigit(runes[i+2]) {
					return "", fmt.Errorf("invalid escape sequence: \\0 followed by digit")
				}
				result.WriteByte(0)
				i++
			case '\n':
				// Line continuation: backslash + LF — both are removed
				i++
			case '\r':
				// Line continuation: backslash + CR (optionally followed by LF)
				i++ // now i points to the \r position; next iteration i++ skips past it
				if i+1 < length && runes[i+1] == '\n' {
					i++ // also skip the \n in CRLF
				}
			case '\u2028', '\u2029':
				// Line continuation: backslash + LS or PS
				i++
			case 'u', 'U':
				// Handle \uXXXX, \UXXXXXXXX, \u{0xXXXX}, \U{0xXXXX}
				parsed, advance, err := parseUnicodeEscape(runes, i)
				if err != nil {
					return "", err
				}
				result.WriteRune(parsed)
				i = advance
			default:
				// Unknown escape: preserve backslash and the character
				result.WriteRune(ch)
				result.WriteRune(nextCh)
				i++
			}
		} else {
			result.WriteRune(ch)
		}
	}

	return result.String(), nil
}

// parseUnicodeEscape handles \uXXXX, \UXXXXXXXX, \u{...}, \U{...} sequences.
// runes[pos] is '\', runes[pos+1] is 'u' or 'U'.
// Returns the decoded rune, the new index (pointing at the last consumed rune), and any error.
func parseUnicodeEscape(runes []rune, pos int) (rune, int, error) {
	marker := runes[pos+1] // 'u' or 'U'
	i := pos

	if i+2 < len(runes) && runes[i+2] == '{' {
		// Handle \u{...} or \U{...}
		i += 3 // skip \, u/U, {
		start := i
		for i < len(runes) && runes[i] != '}' {
			i++
		}
		if i >= len(runes) {
			return 0, 0, fmt.Errorf("invalid Unicode escape: incomplete \\u or \\U sequence")
		}
		hex := string(runes[start:i])
		if strings.HasPrefix(hex, "0x") {
			hex = hex[2:]
		}
		codePoint, err := strconv.ParseInt(hex, 16, 32)
		if err != nil || !utf8.ValidRune(rune(codePoint)) {
			return 0, 0, fmt.Errorf("invalid Unicode escape: \\%c{%s}", marker, hex)
		}
		return rune(codePoint), i, nil
	}

	if marker == 'u' && i+5 < len(runes) {
		// Handle \uXXXX
		hex := string(runes[i+2 : i+6])
		codePoint, err := strconv.ParseInt(hex, 16, 32)
		if err != nil || !utf8.ValidRune(rune(codePoint)) {
			return 0, 0, fmt.Errorf("invalid Unicode escape: \\u%s", hex)
		}
		return rune(codePoint), i + 5, nil
	}

	if marker == 'U' && i+9 < len(runes) {
		// Handle \UXXXXXXXX
		hex := string(runes[i+2 : i+10])
		codePoint, err := strconv.ParseInt(hex, 16, 32)
		if err != nil || !utf8.ValidRune(rune(codePoint)) {
			return 0, 0, fmt.Errorf("invalid Unicode escape: \\U%s", hex)
		}
		return rune(codePoint), i + 9, nil
	}

	return 0, 0, fmt.Errorf("invalid Unicode escape")
}

// consumeExponent advances past an optional exponent part (e/E followed by optional sign and digits).
// Returns the new byte offset.
func consumeExponent(input string, i int) int {
	length := len(input)
	if i < length && (input[i] == 'e' || input[i] == 'E') {
		i++
		if i < length && (input[i] == '+' || input[i] == '-') {
			i++
		}
		for i < length && isDigit(rune(input[i])) {
			i++
		}
	}
	return i
}

// Tokenize tokenizes a JSON5 string with support for unquoted keys,
// hex numbers, escape sequences, comments, and all JSON5 number literals.
// It iterates over the input using proper rune decoding so that multi-byte
// Unicode whitespace characters (e.g. \u00A0, \u2028) are handled correctly.
func Tokenize(input string) []Token {
	var tokens []Token
	length := len(input)
	i := 0

	for i < length {
		// Decode the current rune and its byte width
		ch, size := utf8.DecodeRuneInString(input[i:])

		// Skip whitespace (works for multi-byte whitespace like \u00A0, \u2028)
		if isWhitespace(ch) {
			i += size
			continue
		}

		// Handle single-line (//) and multi-line (/* */) comments
		if ch == '/' && i+size < length {
			next, _ := utf8.DecodeRuneInString(input[i+size:])
			if next == '/' {
				// Single-line comment: advance to end of line
				start := i
				i += 2
				for i < length {
					r, sz := utf8.DecodeRuneInString(input[i:])
					if isLineTerminator(r) {
						break
					}
					i += sz
				}
				tokens = append(tokens, Token{Type: TOKEN_COMMENT, Value: input[start:i]})
				continue
			} else if next == '*' {
				// Multi-line comment
				start := i
				i += 2
				for i < length-1 && !(input[i] == '*' && input[i+1] == '/') {
					i++
				}
				if i < length-1 {
					i += 2 // Skip over the closing */
				} else {
					i = length // unterminated comment — consume rest of input
				}
				tokens = append(tokens, Token{Type: TOKEN_COMMENT, Value: input[start:i]})
				continue
			}
		}

		switch ch {
		case '{':
			tokens = append(tokens, Token{Type: TOKEN_LBRACE, Value: "{"})
			i += size
		case '}':
			tokens = append(tokens, Token{Type: TOKEN_RBRACE, Value: "}"})
			i += size
		case '[':
			tokens = append(tokens, Token{Type: TOKEN_LBRACKET, Value: "["})
			i += size
		case ']':
			tokens = append(tokens, Token{Type: TOKEN_RBRACKET, Value: "]"})
			i += size
		case ':':
			tokens = append(tokens, Token{Type: TOKEN_COLON, Value: ":"})
			i += size
		case ',':
			tokens = append(tokens, Token{Type: TOKEN_COMMA, Value: ","})
			i += size
		case '"', '\'':
			// Quoted string — scan for matching closing quote, handling escapes.
			// We use rune decoding so multi-byte characters don't create false matches.
			start := i
			quote := ch
			i += size // skip opening quote
			for i < length {
				r, sz := utf8.DecodeRuneInString(input[i:])
				if r == '\\' && i+sz < length {
					// Escape sequence: skip the backslash and the next rune
					_, sz2 := utf8.DecodeRuneInString(input[i+sz:])
					i += sz + sz2
				} else if r == quote {
					break
				} else {
					i += sz
				}
			}
			quoteLen := utf8.RuneLen(quote)
			i += quoteLen // skip closing quote
			rawString := input[start+quoteLen : i-quoteLen]
			processedString, err := processEscapeSequences(rawString)
			if err != nil {
				tokens = append(tokens, Token{Type: TOKEN_UNKNOWN, Value: err.Error()})
			} else {
				tokens = append(tokens, Token{Type: TOKEN_STRING, Value: processedString})
			}
		default:
			if ch == '.' && i+1 < length && isDigit(rune(input[i+1])) {
				// Leading decimal point number: .5, .123
				start := i
				i++ // skip the dot
				for i < length && isDigit(rune(input[i])) {
					i++
				}
				i = consumeExponent(input, i)
				tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
			} else if isDigit(ch) || ch == '-' || ch == '+' {
				// Number token including hex, signs, Infinity
				start := i
				if ch == '-' || ch == '+' {
					i++ // skip the sign
				}
				// After sign: check we have more input with a valid number start
				if i >= length || (!isDigit(rune(input[i])) && input[i] != '.' &&
					!strings.HasPrefix(input[i:], "Infinity") &&
					!strings.HasPrefix(input[i:], "0x") && !strings.HasPrefix(input[i:], "0X")) {
					// Bare sign with no valid number following — emit as unknown
					tokens = append(tokens, Token{Type: TOKEN_UNKNOWN, Value: input[start:i]})
				} else if strings.HasPrefix(input[i:], "Infinity") {
					i += 8
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				} else if i+1 < length && (input[i:i+2] == "0x" || input[i:i+2] == "0X") {
					// Hexadecimal number (with optional sign)
					i += 2
					for i < length && isHexDigit(rune(input[i])) {
						i++
					}
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				} else if i < length && input[i] == '.' && i+1 < length && isDigit(rune(input[i+1])) {
					// Signed leading decimal point: +.5, -.5
					i++ // skip dot
					for i < length && isDigit(rune(input[i])) {
						i++
					}
					i = consumeExponent(input, i)
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				} else {
					// Decimal number (handles trailing dot like 5. and exponent signs like 1e+5)
					for i < length && isDigit(rune(input[i])) {
						i++
					}
					// Optional decimal part (including trailing dot)
					if i < length && input[i] == '.' {
						i++
						for i < length && isDigit(rune(input[i])) {
							i++
						}
					}
					i = consumeExponent(input, i)
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				}
			} else if isIdentifierStart(ch) {
				// Unquoted key or identifier (rune-safe iteration)
				start := i
				i += size
				for i < length {
					r, sz := utf8.DecodeRuneInString(input[i:])
					if !isIdentifierPart(r) {
						break
					}
					i += sz
				}
				unquotedString := input[start:i]
				switch unquotedString {
				case "true":
					tokens = append(tokens, Token{Type: TOKEN_TRUE, Value: "true"})
				case "false":
					tokens = append(tokens, Token{Type: TOKEN_FALSE, Value: "false"})
				case "null":
					tokens = append(tokens, Token{Type: TOKEN_NULL, Value: "null"})
				case "Infinity":
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: "Infinity"})
				case "NaN":
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: "NaN"})
				default:
					tokens = append(tokens, Token{Type: TOKEN_STRING, Value: unquotedString})
				}
			} else {
				// Unknown token
				tokens = append(tokens, Token{Type: TOKEN_UNKNOWN, Value: string(ch)})
				i += size
			}
		}
	}

	return tokens
}
