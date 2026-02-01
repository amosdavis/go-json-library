package json

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Marshal converts a Go value to a JSON string
func Marshal(v interface{}) ([]byte, error) {
	var sb strings.Builder
	if err := marshalValue(&sb, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}

// MarshalIndent converts a Go value to a formatted JSON string
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	data, err := Marshal(v)
	if err != nil {
		return nil, err
	}
	return formatJSON(data, prefix, indent)
}

func marshalValue(sb *strings.Builder, v reflect.Value) error {
	if !v.IsValid() {
		sb.WriteString("null")
		return nil
	}

	// Handle interface and pointer types
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			sb.WriteString("null")
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsInf(f, 0) || math.IsNaN(f) {
			return fmt.Errorf("json: unsupported value: %v", f)
		}
		sb.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
	case reflect.String:
		marshalString(sb, v.String())
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			sb.WriteString("null")
			return nil
		}
		sb.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				sb.WriteByte(',')
			}
			if err := marshalValue(sb, v.Index(i)); err != nil {
				return err
			}
		}
		sb.WriteByte(']')
	case reflect.Map:
		if v.IsNil() {
			sb.WriteString("null")
			return nil
		}
		sb.WriteByte('{')
		keys := v.MapKeys()
		// Sort keys for deterministic output
		sortedKeys := make([]string, 0, len(keys))
		keyMap := make(map[string]reflect.Value)
		for _, k := range keys {
			if k.Kind() != reflect.String {
				return fmt.Errorf("json: map key must be string, got %v", k.Kind())
			}
			keyStr := k.String()
			sortedKeys = append(sortedKeys, keyStr)
			keyMap[keyStr] = k
		}
		sort.Strings(sortedKeys)
		for i, keyStr := range sortedKeys {
			if i > 0 {
				sb.WriteByte(',')
			}
			marshalString(sb, keyStr)
			sb.WriteByte(':')
			if err := marshalValue(sb, v.MapIndex(keyMap[keyStr])); err != nil {
				return err
			}
		}
		sb.WriteByte('}')
	case reflect.Struct:
		return marshalStruct(sb, v)
	default:
		return fmt.Errorf("json: unsupported type: %v", v.Kind())
	}
	return nil
}

func marshalStruct(sb *strings.Builder, v reflect.Value) error {
	sb.WriteByte('{')
	t := v.Type()
	first := true
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}

		fieldValue := v.Field(i)
		name := field.Name
		omitempty := false

		// Handle json struct tags
		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if parts[0] != "" {
				name = parts[0]
			}
			for _, opt := range parts[1:] {
				if opt == "omitempty" {
					omitempty = true
				}
			}
		}

		if omitempty && isEmptyValue(fieldValue) {
			continue
		}

		if !first {
			sb.WriteByte(',')
		}
		first = false

		marshalString(sb, name)
		sb.WriteByte(':')
		if err := marshalValue(sb, fieldValue); err != nil {
			return err
		}
	}
	sb.WriteByte('}')
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func marshalString(sb *strings.Builder, s string) {
	sb.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString("\\\"")
		case '\\':
			sb.WriteString("\\\\")
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			if r < 0x20 {
				sb.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				sb.WriteRune(r)
			}
		}
	}
	sb.WriteByte('"')
}

func formatJSON(data []byte, prefix, indent string) ([]byte, error) {
	var sb strings.Builder
	level := 0
	inString := false
	escape := false

	writeIndent := func() {
		sb.WriteString(prefix)
		for i := 0; i < level; i++ {
			sb.WriteString(indent)
		}
	}

	for i := 0; i < len(data); i++ {
		c := data[i]

		if escape {
			sb.WriteByte(c)
			escape = false
			continue
		}

		if c == '\\' && inString {
			sb.WriteByte(c)
			escape = true
			continue
		}

		if c == '"' {
			inString = !inString
			sb.WriteByte(c)
			continue
		}

		if inString {
			sb.WriteByte(c)
			continue
		}

		switch c {
		case '{', '[':
			sb.WriteByte(c)
			level++
			// Check if empty
			next := i + 1
			for next < len(data) && (data[next] == ' ' || data[next] == '\t' || data[next] == '\n' || data[next] == '\r') {
				next++
			}
			if next < len(data) && ((c == '{' && data[next] == '}') || (c == '[' && data[next] == ']')) {
				// Empty object/array
				level--
			} else {
				sb.WriteByte('\n')
				writeIndent()
			}
		case '}', ']':
			level--
			// Check if preceded by content
			if i > 0 && data[i-1] != '{' && data[i-1] != '[' {
				sb.WriteByte('\n')
				writeIndent()
			}
			sb.WriteByte(c)
		case ',':
			sb.WriteByte(c)
			sb.WriteByte('\n')
			writeIndent()
		case ':':
			sb.WriteByte(c)
			sb.WriteByte(' ')
		default:
			sb.WriteByte(c)
		}
	}

	return []byte(sb.String()), nil
}
