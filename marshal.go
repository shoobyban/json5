package json5

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshal converts an interface{} into a JSON5 string.
func Marshal(value interface{}) (string, error) {
	return marshalValue(value, "", 0)
}

// MarshalIndent converts an interface{} into a JSON5 string with indentation.
func MarshalIndent(value interface{}, indent string) (string, error) {
	return marshalValue(value, indent, 0)
}

// marshalValue recursively converts a Go value into a JSON5 string
func marshalValue(value interface{}, indent string, depth int) (string, error) {
	switch v := value.(type) {
	case nil:
		return "null", nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	case string:
		return marshalString(v), nil
	case []interface{}:
		return marshalArray(v, indent, depth)
	case map[string]interface{}:
		return marshalObject(v, indent, depth)
	default:
		// Handle other types if needed (custom types, etc.)
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// marshalString handles string values and escapes necessary characters
func marshalString(s string) string {
	// Escape necessary characters in the string
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"\"", "\\\"",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	)
	escaped := replacer.Replace(s)
	return fmt.Sprintf("\"%s\"", escaped)
}

// marshalArray handles slices of interface{} and converts them into JSON5 arrays
func marshalArray(array []interface{}, indent string, depth int) (string, error) {
	var sb strings.Builder
	sb.WriteString("[")

	newIndent := strings.Repeat(indent, depth+1)
	first := true

	for _, item := range array {
		if !first {
			sb.WriteString(", ")
		}
		first = false

		itemStr, err := marshalValue(item, indent, depth+1)
		if err != nil {
			return "", err
		}

		sb.WriteString("\n")
		sb.WriteString(newIndent)
		sb.WriteString(itemStr)
	}

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(indent, depth))
	sb.WriteString("]")
	return sb.String(), nil
}

// marshalObject handles maps and converts them into JSON5 objects
func marshalObject(obj map[string]interface{}, indent string, depth int) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")

	newIndent := strings.Repeat(indent, depth+1)

	for key, value := range obj {
		keyStr := marshalKey(key)

		valueStr, err := marshalValue(value, indent, depth+1)
		if err != nil {
			return "", err
		}

		sb.WriteString("\n")
		sb.WriteString(newIndent)
		sb.WriteString(keyStr)
		sb.WriteString(": ")
		sb.WriteString(valueStr)
		sb.WriteString(",")
	}

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat(indent, depth))
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
