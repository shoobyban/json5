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

// isWhitespace checks if a character is a whitespace character
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
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

// processEscapeSequences converts escape sequences such as \n, \t, \uXXXX, \UXXXXXXXX, \u{0x1FA}, and \U{0x1FA} into their actual representations
func processEscapeSequences(input string) (string, error) {
	var result strings.Builder
	length := len(input)

	for i := 0; i < length; i++ {
		ch := input[i]

		if ch == '\\' && i+1 < length {
			nextCh := input[i+1]
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
			case 'u', 'U': // Handle \uXXXX, \UXXXXXXXX, \u{0xXXXX}, \U{0xXXXX}
				if i+2 < length && input[i+2] == '{' {
					// Handle \u{...} or \U{...}
					i += 3
					start := i
					for i < length && input[i] != '}' {
						i++
					}
					if i >= length {
						return "", fmt.Errorf("invalid Unicode escape: incomplete \\u or \\U sequence")
					}
					hex := input[start:i]
					if strings.HasPrefix(hex, "0x") {
						hex = hex[2:] // Remove 0x if present
					}
					codePoint, err := strconv.ParseInt(hex, 16, 32)
					if err != nil || !utf8.ValidRune(rune(codePoint)) {
						return "", fmt.Errorf("invalid Unicode escape: \\u{%s}", hex)
					}
					result.WriteRune(rune(codePoint))
				} else if nextCh == 'u' && i+5 < length {
					// Handle \uXXXX
					hex := input[i+2 : i+6]
					codePoint, err := strconv.ParseInt(hex, 16, 32)
					if err != nil || !utf8.ValidRune(rune(codePoint)) {
						return "", fmt.Errorf("invalid Unicode escape: \\u%s", hex)
					}
					result.WriteRune(rune(codePoint))
					i += 5 // Move past the 4 hex digits
				} else if nextCh == 'U' && i+9 < length {
					// Handle \UXXXXXXXX
					hex := input[i+2 : i+10]
					codePoint, err := strconv.ParseInt(hex, 16, 32)
					if err != nil || !utf8.ValidRune(rune(codePoint)) {
						return "", fmt.Errorf("invalid Unicode escape: \\U%s", hex)
					}
					result.WriteRune(rune(codePoint))
					i += 9 // Move past the 8 hex digits
				} else {
					return "", fmt.Errorf("invalid Unicode escape")
				}
			default:
				result.WriteByte(ch)
			}
		} else {
			result.WriteByte(ch)
		}
	}

	return result.String(), nil
}

// tokenize is a function to tokenize a JSON5 string with unquoted key and hex number support, and escape sequence handling
func Tokenize(input string) []Token {
	var tokens []Token
	length := len(input)
	i := 0

	for i < length {
		ch := rune(input[i])

		// Skip whitespaces
		if isWhitespace(ch) {
			i++
			continue
		}

		// Handle single-line (//) and multi-line (/* */) comments
		if ch == '/' && i+1 < length {
			nextCh := rune(input[i+1])
			if nextCh == '/' {
				// Single-line comment
				start := i
				i += 2
				for i < length && input[i] != '\n' {
					i++
				}
				tokens = append(tokens, Token{Type: TOKEN_COMMENT, Value: input[start:i]})
				continue
			} else if nextCh == '*' {
				// Multi-line comment
				start := i
				i += 2
				for i < length-1 && !(input[i] == '*' && input[i+1] == '/') {
					i++
				}
				i += 2 // Skip over the closing */
				tokens = append(tokens, Token{Type: TOKEN_COMMENT, Value: input[start:i]})
				continue
			}
		}

		switch ch {
		case '{':
			tokens = append(tokens, Token{Type: TOKEN_LBRACE, Value: "{"})
			i++
		case '}':
			tokens = append(tokens, Token{Type: TOKEN_RBRACE, Value: "}"})
			i++
		case '[':
			tokens = append(tokens, Token{Type: TOKEN_LBRACKET, Value: "["})
			i++
		case ']':
			tokens = append(tokens, Token{Type: TOKEN_RBRACKET, Value: "]"})
			i++
		case ':':
			tokens = append(tokens, Token{Type: TOKEN_COLON, Value: ":"})
			i++
		case ',':
			tokens = append(tokens, Token{Type: TOKEN_COMMA, Value: ","})
			i++
		case '"', '\'':
			// String token (quoted), handle both keys and values
			start := i
			quote := ch
			i++
			for i < length && rune(input[i]) != quote {
				if input[i] == '\\' && i+1 < length {
					// Escape sequence detected, skip it
					i += 2
				} else {
					i++
				}
			}
			i++                                                       // Skip the closing quote
			rawString := input[start+1 : i-1]                         // Remove quotes
			processedString, err := processEscapeSequences(rawString) // Handle escape sequences
			if err != nil {
				tokens = append(tokens, Token{Type: TOKEN_UNKNOWN, Value: err.Error()})
			} else {
				tokens = append(tokens, Token{Type: TOKEN_STRING, Value: processedString})
			}
		default:
			if isDigit(ch) || ch == '-' {
				// Number token, including support for hex
				start := i
				if input[i:i+2] == "0x" || input[i:i+2] == "0X" {
					// Hexadecimal number
					i += 2
					for i < length && isHexDigit(rune(input[i])) {
						i++
					}
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				} else {
					// Decimal number
					for i < length && (isDigit(rune(input[i])) || input[i] == '.' || input[i] == 'e' || input[i] == 'E') {
						i++
					}
					tokens = append(tokens, Token{Type: TOKEN_NUMBER, Value: input[start:i]})
				}
			} else if isIdentifierStart(ch) {
				// Unquoted key or identifier
				start := i
				for i < length && isIdentifierPart(rune(input[i])) {
					i++
				}
				unquotedString := input[start:i]
				// Check if it's a boolean or null literal
				switch unquotedString {
				case "true":
					tokens = append(tokens, Token{Type: TOKEN_TRUE, Value: "true"})
				case "false":
					tokens = append(tokens, Token{Type: TOKEN_FALSE, Value: "false"})
				case "null":
					tokens = append(tokens, Token{Type: TOKEN_NULL, Value: "null"})
				default:
					tokens = append(tokens, Token{Type: TOKEN_STRING, Value: unquotedString})
				}
			} else {
				// Unknown token (for simplicity)
				tokens = append(tokens, Token{Type: TOKEN_UNKNOWN, Value: string(ch)})
				i++
			}
		}
	}

	return tokens
}
