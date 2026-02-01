package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal parses JSON data and stores the result in the value pointed to by v
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("json: Unmarshal requires non-nil pointer")
	}

	parser := NewParser(string(data))
	parsed, err := parser.Parse()
	if err != nil {
		return err
	}

	return unmarshalValue(parsed, rv.Elem())
}

func unmarshalValue(data interface{}, v reflect.Value) error {
	if data == nil {
		// Set to zero value for nil
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return unmarshalValue(data, v.Elem())
	}

	// Handle interface{} - just set directly
	if v.Kind() == reflect.Interface {
		v.Set(reflect.ValueOf(data))
		return nil
	}

	switch d := data.(type) {
	case bool:
		return unmarshalBool(d, v)
	case float64:
		return unmarshalNumber(d, v)
	case string:
		return unmarshalString(d, v)
	case []interface{}:
		return unmarshalArray(d, v)
	case map[string]interface{}:
		return unmarshalObject(d, v)
	default:
		return fmt.Errorf("json: unknown type %T", data)
	}
}

func unmarshalBool(data bool, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(data)
	default:
		return fmt.Errorf("json: cannot unmarshal bool into %v", v.Type())
	}
	return nil
}

func unmarshalNumber(data float64, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(data))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(data))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(data)
	default:
		return fmt.Errorf("json: cannot unmarshal number into %v", v.Type())
	}
	return nil
}

func unmarshalString(data string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(data)
	default:
		return fmt.Errorf("json: cannot unmarshal string into %v", v.Type())
	}
	return nil
}

func unmarshalArray(data []interface{}, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(v.Type(), len(data), len(data))
		for i, item := range data {
			if err := unmarshalValue(item, slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)
	case reflect.Array:
		if v.Len() < len(data) {
			return fmt.Errorf("json: array overflow")
		}
		for i, item := range data {
			if err := unmarshalValue(item, v.Index(i)); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("json: cannot unmarshal array into %v", v.Type())
	}
	return nil
}

func unmarshalObject(data map[string]interface{}, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("json: map key must be string")
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
		for key, val := range data {
			elemType := v.Type().Elem()
			elem := reflect.New(elemType).Elem()
			if err := unmarshalValue(val, elem); err != nil {
				return err
			}
			v.SetMapIndex(reflect.ValueOf(key), elem)
		}
	case reflect.Struct:
		return unmarshalStruct(data, v)
	default:
		return fmt.Errorf("json: cannot unmarshal object into %v", v.Type())
	}
	return nil
}

func unmarshalStruct(data map[string]interface{}, v reflect.Value) error {
	t := v.Type()

	// Build field name map
	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		name := field.Name
		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				continue
			}
			if parts[0] != "" {
				name = parts[0]
			}
		}
		fieldMap[name] = i
		fieldMap[strings.ToLower(name)] = i // case-insensitive fallback
	}

	for key, val := range data {
		idx, ok := fieldMap[key]
		if !ok {
			idx, ok = fieldMap[strings.ToLower(key)]
		}
		if !ok {
			continue // ignore unknown fields
		}
		field := v.Field(idx)
		if !field.CanSet() {
			continue
		}
		if err := unmarshalValue(val, field); err != nil {
			return fmt.Errorf("json: error unmarshaling field %s: %v", key, err)
		}
	}
	return nil
}

// Valid reports whether data is a valid JSON encoding
func Valid(data []byte) bool {
	parser := NewParser(string(data))
	_, err := parser.Parse()
	return err == nil
}

// Number represents a JSON number literal
type Number string

// String returns the literal text of the number
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}
