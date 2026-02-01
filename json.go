// Package json provides a native Go implementation of JSON encoding and decoding
// based on RFC 8259 (The JavaScript Object Notation Data Interchange Format).
//
// This package provides Marshal and Unmarshal functions similar to the standard
// library's encoding/json package, but implemented from scratch without any
// external dependencies.
//
// Basic usage:
//
//	// Marshal a Go value to JSON
//	data, err := json.Marshal(myStruct)
//
//	// Unmarshal JSON to a Go value
//	var result MyStruct
//	err := json.Unmarshal(data, &result)
//
// The package supports:
//   - All JSON value types: objects, arrays, strings, numbers, booleans, null
//   - Struct tags for field name mapping and omitempty
//   - Unicode escape sequences
//   - Pretty printing with MarshalIndent
package json

// Encoder writes JSON values to an output stream (future enhancement)
// Decoder reads and decodes JSON values from an input stream (future enhancement)
