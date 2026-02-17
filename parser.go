package json5

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func UnMarshal(json5 string) (any, error) {
	allTokens := Tokenize(json5)

	// Filter out comment tokens so the parser handles JSON5 comments
	tokens := make([]Token, 0, len(allTokens))
	for _, t := range allTokens {
		if t.Type != TOKEN_COMMENT {
			tokens = append(tokens, t)
		}
	}

	tokenLen := len(tokens)

	if tokenLen == 0 {
		return nil, nil
	}

	i := 1 // start from the second token (skip the first one we already checked)
	if tokens[0].Type == TOKEN_LBRACE {
		obj, err := parseObject(tokens, &i, tokenLen)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}

	if tokens[0].Type == TOKEN_LBRACKET {
		arr, err := parseArray(tokens, &i, tokenLen)
		if err != nil {
			return nil, err
		}
		return arr, nil
	}

	if tokens[0].Type == TOKEN_STRING {
		return tokens[0].Value, nil
	}

	if tokens[0].Type == TOKEN_NUMBER {
		return parseNumber(tokens[0].Value)
	}

	if tokens[0].Type == TOKEN_TRUE {
		return true, nil
	}

	if tokens[0].Type == TOKEN_FALSE {
		return false, nil
	}

	if tokens[0].Type == TOKEN_NULL {
		return nil, nil
	}

	return nil, fmt.Errorf("expected '{', '[', number, null or boolean but found '%s'", tokens[0].Value)
}

// Unmarshal is an alias for UnMarshal that follows Go naming conventions.
func Unmarshal(json5 string) (any, error) {
	return UnMarshal(json5)
}

// parseObject parses the tokens as a JSON5 object and returns a map[string]interface{}
func parseObject(tokens []Token, i *int, tokenLen int) (map[string]any, error) {
	result := make(map[string]any)

	for *i < tokenLen {
		// If we encounter a closing brace, we're done with the object
		if tokens[*i].Type == TOKEN_RBRACE {
			*i++ // Move past the closing brace
			break
		}

		// Parse the key (it should be a string or unquoted identifier)
		keyToken := tokens[*i]
		if keyToken.Type != TOKEN_STRING {
			return nil, fmt.Errorf("expected a string for key but found '%s'", keyToken.Value)
		}
		key := keyToken.Value
		*i++

		// Expect a colon after the key
		if *i >= tokenLen {
			return nil, fmt.Errorf("unexpected end of input: expected ':' after key '%s'", key)
		}
		if tokens[*i].Type != TOKEN_COLON {
			return nil, fmt.Errorf("expected ':' after key '%s' but found '%s'", key, tokens[*i].Value)
		}
		*i++

		// Parse the value for the key
		value, err := parseValue(tokens, i, tokenLen)
		if err != nil {
			return nil, err
		}

		// Add the key-value pair to the result
		result[key] = value

		// After the value, we should either find a comma or a closing brace
		if *i >= tokenLen {
			return nil, fmt.Errorf("unexpected end of input: expected ',' or '}'")
		}
		if tokens[*i].Type == TOKEN_COMMA {
			*i++ // Move past the comma
		} else if tokens[*i].Type == TOKEN_RBRACE {
			*i++ // Move past the closing brace
			break
		} else {
			return nil, fmt.Errorf("expected ',' or '}' but found '%s'", tokens[*i].Value)
		}
	}

	return result, nil
}

// parseArray parses the tokens as a JSON5 array and returns a []interface{}
func parseArray(tokens []Token, i *int, tokenLen int) ([]any, error) {
	var result []any

	for *i < tokenLen {
		// If we encounter a closing bracket, we're done with the array
		if tokens[*i].Type == TOKEN_RBRACKET {
			*i++ // Move past the closing bracket
			break
		}

		// Parse the next value
		value, err := parseValue(tokens, i, tokenLen)
		if err != nil {
			return nil, err
		}
		result = append(result, value)

		// After the value, we should either find a comma or a closing bracket
		if *i >= tokenLen {
			return nil, fmt.Errorf("unexpected end of input: expected ',' or ']'")
		}
		if tokens[*i].Type == TOKEN_COMMA {
			*i++ // Move past the comma
		} else if tokens[*i].Type == TOKEN_RBRACKET {
			*i++ // Move past the closing bracket
			break
		} else {
			return nil, fmt.Errorf("expected ',' or ']' but found '%s'", tokens[*i].Value)
		}
	}

	return result, nil
}

// parseValue parses a value (string, number, boolean, null, object, or array)
func parseValue(tokens []Token, i *int, tokenLen int) (any, error) {
	if *i >= tokenLen {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch tokens[*i].Type {
	case TOKEN_STRING:
		value := tokens[*i].Value
		*i++
		return value, nil
	case TOKEN_NUMBER:
		numberStr := tokens[*i].Value
		*i++
		return parseNumber(numberStr)
	case TOKEN_TRUE:
		*i++
		return true, nil
	case TOKEN_FALSE:
		*i++
		return false, nil
	case TOKEN_NULL:
		*i++
		return nil, nil // Return nil when the value is the literal `null`
	case TOKEN_LBRACE:
		*i++ // Move past the opening brace
		return parseObject(tokens, i, tokenLen)
	case TOKEN_LBRACKET:
		*i++ // Move past the opening bracket
		return parseArray(tokens, i, tokenLen)
	default:
		return nil, fmt.Errorf("unexpected token: '%s'", tokens[*i].Value)
	}
}

// parseNumber parses a JSON5 number string into its Go representation.
// Handles integers, floats, hexadecimal (with optional sign), Infinity,
// -Infinity, +Infinity, NaN, leading decimal point (.5), trailing
// decimal point (5.), and exponent signs (1e+5).
func parseNumber(numberStr string) (any, error) {
	// Handle Infinity and NaN
	switch numberStr {
	case "Infinity", "+Infinity":
		return math.Inf(1), nil
	case "-Infinity":
		return math.Inf(-1), nil
	case "NaN":
		return math.NaN(), nil
	}

	// Determine sign and strip it for parsing.
	// Go's strconv.Atoi/ParseFloat accept '-' but not '+'.
	parseable := numberStr
	sign := int64(1)
	if len(parseable) > 0 && (parseable[0] == '+' || parseable[0] == '-') {
		if parseable[0] == '-' {
			sign = -1
		}
		parseable = parseable[1:]
	}

	// Check if the number is hexadecimal
	if strings.HasPrefix(parseable, "0x") || strings.HasPrefix(parseable, "0X") {
		num, err := strconv.ParseInt(parseable, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid hexadecimal number: '%s'", numberStr)
		}
		num *= sign
		if num >= math.MinInt && num <= math.MaxInt {
			return int(num), nil
		}
		return num, nil
	}

	// Re-add '-' for negative decimals so strconv handles them natively
	if sign == -1 {
		parseable = "-" + parseable
	}

	// Parse as a regular decimal number (int or float)
	if num, err := strconv.Atoi(parseable); err == nil {
		return num, nil
	}
	if num, err := strconv.ParseFloat(parseable, 64); err == nil {
		return num, nil
	}
	return nil, fmt.Errorf("invalid number: '%s'", numberStr)
}
