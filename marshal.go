package json5

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Marshal converts an interface{} into a JSON5 string.
func Marshal(value any) (string, error) {
	return marshalValue(value, "", 0)
}

// MarshalIndent converts an interface{} into a JSON5 string with indentation.
func MarshalIndent(value any, indent string) (string, error) {
	return marshalValue(value, indent, 0)
}

// marshalValue recursively converts a Go value into a JSON5 string
func marshalValue(value any, indent string, depth int) (string, error) {
	switch v := value.(type) {
	case nil:
		return "null", nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", v), nil
	case float32:
		f := float64(v)
		if math.IsInf(f, 1) {
			return "Infinity", nil
		} else if math.IsInf(f, -1) {
			return "-Infinity", nil
		} else if math.IsNaN(f) {
			return "NaN", nil
		}
		return fmt.Sprintf("%v", v), nil
	case float64:
		if math.IsInf(v, 1) {
			return "Infinity", nil
		} else if math.IsInf(v, -1) {
			return "-Infinity", nil
		} else if math.IsNaN(v) {
			return "NaN", nil
		}
		return fmt.Sprintf("%v", v), nil
	case string:
		return marshalString(v), nil
	case []any:
		return marshalArray(v, indent, depth)
	case map[string]any:
		return marshalObject(v, indent, depth)
	default:
		// Handle other types if needed (custom types, etc.)
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// marshalString handles string values and escapes necessary characters.
// Escapes backslash, double quote, and all JSON5 control characters:
// \n, \r, \t, \b, \f, \v, and \0 (null byte).
func marshalString(s string) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\v':
			sb.WriteString("\\v")
		case 0:
			sb.WriteString("\\0")
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

// marshalArray handles slices of interface{} and converts them into JSON5 arrays
func marshalArray(array []any, indent string, depth int) (string, error) {
	if len(array) == 0 {
		return "[]", nil
	}

	var sb strings.Builder
	compact := indent == ""
	sb.WriteString("[")

	for i, item := range array {
		itemStr, err := marshalValue(item, indent, depth+1)
		if err != nil {
			return "", err
		}

		if compact {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(itemStr)
		} else {
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat(indent, depth+1))
			sb.WriteString(itemStr)
			sb.WriteString(",")
		}
	}

	if !compact {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(indent, depth))
	}
	sb.WriteString("]")
	return sb.String(), nil
}

// marshalObject handles maps and converts them into JSON5 objects
func marshalObject(obj map[string]any, indent string, depth int) (string, error) {
	if len(obj) == 0 {
		return "{}", nil
	}

	var sb strings.Builder
	compact := indent == ""
	sb.WriteString("{")

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		value := obj[key]
		keyStr := marshalKey(key)

		valueStr, err := marshalValue(value, indent, depth+1)
		if err != nil {
			return "", err
		}

		if compact {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(keyStr)
			sb.WriteString(": ")
			sb.WriteString(valueStr)
		} else {
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat(indent, depth+1))
			sb.WriteString(keyStr)
			sb.WriteString(": ")
			sb.WriteString(valueStr)
			sb.WriteString(",")
		}
	}

	if !compact {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(indent, depth))
	}
	sb.WriteString("}")
	return sb.String(), nil
}

// marshalKey checks if a key can be unquoted in JSON5 (simple identifier) or must be quoted
func marshalKey(key string) string {
	// Check if the key can be unquoted (simple identifier rules for JSON5)
	if isSimpleIdentifier(key) {
		return key
	}
	// Otherwise, quote the key
	return strconv.Quote(key)
}

// isSimpleIdentifier checks if a string qualifies as a simple identifier (unquoted in JSON5)
func isSimpleIdentifier(key string) bool {
	if len(key) == 0 {
		return false
	}
	for i, ch := range key {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '$' || (i > 0 && ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}
