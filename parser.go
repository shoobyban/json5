package json5

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Parse(json5 string) (interface{}, error) {
	tokens := Tokenize(json5)

	tokenLen := len(tokens)

	if tokenLen == 0 {
		return nil, nil
	}

	i := 1 // start from the second token (skip the first one we already checked)
	if tokens[0].Type == TOKEN_LBRACE {
		// Parse the object
		obj, err := parseObject(tokens, &i, tokenLen)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}

	if tokens[0].Type == TOKEN_LBRACKET {
		// Parse the array
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
		numberStr := tokens[0].Value
		if num, err := strconv.Atoi(numberStr); err == nil {
			return num, nil
		} else if num, err := strconv.ParseFloat(numberStr, 64); err == nil {
			return num, nil
		} else {
			return nil, fmt.Errorf("invalid number: '%s'", numberStr)
		}
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

// parseObject parses the tokens as a JSON5 object and returns a map[string]interface{}
func parseObject(tokens []Token, i *int, tokenLen int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

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
func parseArray(tokens []Token, i *int, tokenLen int) ([]interface{}, error) {
	var result []interface{}

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
func parseValue(tokens []Token, i *int, tokenLen int) (interface{}, error) {
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
		// Check if the number is hexadecimal
		if strings.HasPrefix(numberStr, "0x") || strings.HasPrefix(numberStr, "0X") {
			// Parse the hexadecimal number
			num, err := strconv.ParseInt(numberStr, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid hexadecimal number: '%s'", numberStr)
			}
			if abs(num) < math.MaxInt {
				return int(num), nil
			}
			return num, nil
		} else {
			// Parse as a regular decimal number (int or float)
			if num, err := strconv.Atoi(numberStr); err == nil {
				return num, nil
			} else if num, err := strconv.ParseFloat(numberStr, 64); err == nil {
				return num, nil
			} else {
				return nil, fmt.Errorf("invalid number: '%s'", numberStr)
			}
		}
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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
