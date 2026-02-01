package json

import (
	"fmt"
	"strconv"
	"strings"
)

// Value represents a dynamic JSON value that can be navigated without structs
type Value struct {
	data interface{}
}

// Load parses a JSON string and returns a Value for dynamic access
func Load(jsonStr string) (*Value, error) {
	parser := NewParser(jsonStr)
	data, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return &Value{data: data}, nil
}

// LoadBytes parses JSON bytes and returns a Value for dynamic access
func LoadBytes(jsonBytes []byte) (*Value, error) {
	return Load(string(jsonBytes))
}

// Raw returns the underlying Go value (map, slice, string, float64, bool, or nil)
func (v *Value) Raw() interface{} {
	if v == nil {
		return nil
	}
	return v.data
}

// Get retrieves a nested value by key (for objects)
// Supports dot notation: v.Get("user.address.city")
func (v *Value) Get(key string) *Value {
	if v == nil || v.data == nil {
		return &Value{data: nil}
	}

	// Handle dot notation
	if strings.Contains(key, ".") {
		parts := strings.SplitN(key, ".", 2)
		return v.Get(parts[0]).Get(parts[1])
	}

	if obj, ok := v.data.(map[string]interface{}); ok {
		return &Value{data: obj[key]}
	}
	return &Value{data: nil}
}

// Index retrieves a value by index (for arrays)
func (v *Value) Index(i int) *Value {
	if v == nil || v.data == nil {
		return &Value{data: nil}
	}
	if arr, ok := v.data.([]interface{}); ok {
		if i >= 0 && i < len(arr) {
			return &Value{data: arr[i]}
		}
	}
	return &Value{data: nil}
}

// String returns the value as a string
func (v *Value) String() string {
	if v == nil || v.data == nil {
		return ""
	}
	switch val := v.data.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Int returns the value as an int
func (v *Value) Int() int {
	if v == nil || v.data == nil {
		return 0
	}
	if num, ok := v.data.(float64); ok {
		return int(num)
	}
	return 0
}

// Int64 returns the value as an int64
func (v *Value) Int64() int64 {
	if v == nil || v.data == nil {
		return 0
	}
	if num, ok := v.data.(float64); ok {
		return int64(num)
	}
	return 0
}

// Float returns the value as a float64
func (v *Value) Float() float64 {
	if v == nil || v.data == nil {
		return 0
	}
	if num, ok := v.data.(float64); ok {
		return num
	}
	return 0
}

// Bool returns the value as a bool
func (v *Value) Bool() bool {
	if v == nil || v.data == nil {
		return false
	}
	if b, ok := v.data.(bool); ok {
		return b
	}
	return false
}

// IsNull returns true if the value is null
func (v *Value) IsNull() bool {
	return v == nil || v.data == nil
}

// IsObject returns true if the value is a JSON object
func (v *Value) IsObject() bool {
	if v == nil {
		return false
	}
	_, ok := v.data.(map[string]interface{})
	return ok
}

// IsArray returns true if the value is a JSON array
func (v *Value) IsArray() bool {
	if v == nil {
		return false
	}
	_, ok := v.data.([]interface{})
	return ok
}

// IsString returns true if the value is a string
func (v *Value) IsString() bool {
	if v == nil {
		return false
	}
	_, ok := v.data.(string)
	return ok
}

// IsNumber returns true if the value is a number
func (v *Value) IsNumber() bool {
	if v == nil {
		return false
	}
	_, ok := v.data.(float64)
	return ok
}

// IsBool returns true if the value is a boolean
func (v *Value) IsBool() bool {
	if v == nil {
		return false
	}
	_, ok := v.data.(bool)
	return ok
}

// Len returns the length of an array or object, or 0 otherwise
func (v *Value) Len() int {
	if v == nil || v.data == nil {
		return 0
	}
	switch val := v.data.(type) {
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	case string:
		return len(val)
	}
	return 0
}

// Keys returns the keys of an object, or nil if not an object
func (v *Value) Keys() []string {
	if v == nil || v.data == nil {
		return nil
	}
	if obj, ok := v.data.(map[string]interface{}); ok {
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

// Values returns the values of an array, or nil if not an array
func (v *Value) Values() []*Value {
	if v == nil || v.data == nil {
		return nil
	}
	if arr, ok := v.data.([]interface{}); ok {
		values := make([]*Value, len(arr))
		for i, item := range arr {
			values[i] = &Value{data: item}
		}
		return values
	}
	return nil
}

// ForEach iterates over array elements or object key-value pairs
// For arrays: fn receives (index as string, value)
// For objects: fn receives (key, value)
func (v *Value) ForEach(fn func(key string, val *Value)) {
	if v == nil || v.data == nil {
		return
	}
	switch val := v.data.(type) {
	case []interface{}:
		for i, item := range val {
			fn(strconv.Itoa(i), &Value{data: item})
		}
	case map[string]interface{}:
		for k, item := range val {
			fn(k, &Value{data: item})
		}
	}
}

// Map returns the value as a map (for objects)
func (v *Value) Map() map[string]*Value {
	if v == nil || v.data == nil {
		return nil
	}
	if obj, ok := v.data.(map[string]interface{}); ok {
		result := make(map[string]*Value, len(obj))
		for k, val := range obj {
			result[k] = &Value{data: val}
		}
		return result
	}
	return nil
}

// Array returns the value as a slice of Values (for arrays)
func (v *Value) Array() []*Value {
	return v.Values()
}

// StringArray returns the value as a slice of strings (for string arrays)
func (v *Value) StringArray() []string {
	if v == nil || v.data == nil {
		return nil
	}
	if arr, ok := v.data.([]interface{}); ok {
		result := make([]string, len(arr))
		for i, item := range arr {
			if s, ok := item.(string); ok {
				result[i] = s
			} else {
				result[i] = fmt.Sprintf("%v", item)
			}
		}
		return result
	}
	return nil
}

// IntArray returns the value as a slice of ints (for number arrays)
func (v *Value) IntArray() []int {
	if v == nil || v.data == nil {
		return nil
	}
	if arr, ok := v.data.([]interface{}); ok {
		result := make([]int, len(arr))
		for i, item := range arr {
			if num, ok := item.(float64); ok {
				result[i] = int(num)
			}
		}
		return result
	}
	return nil
}

// Exists returns true if the value exists and is not null
func (v *Value) Exists() bool {
	return v != nil && v.data != nil
}

// Pretty returns the value as a pretty-printed JSON string
func (v *Value) Pretty() string {
	if v == nil || v.data == nil {
		return "null"
	}
	data, err := MarshalIndent(v.data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v.data)
	}
	return string(data)
}

// Compact returns the value as a compact JSON string
func (v *Value) Compact() string {
	if v == nil || v.data == nil {
		return "null"
	}
	data, err := Marshal(v.data)
	if err != nil {
		return fmt.Sprintf("%v", v.data)
	}
	return string(data)
}
