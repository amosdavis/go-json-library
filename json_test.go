package json

import (
	"reflect"
	"testing"
)

func TestUnmarshalBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"null", "null", nil},
		{"true", "true", true},
		{"false", "false", false},
		{"integer", "42", float64(42)},
		{"negative", "-123", float64(-123)},
		{"float", "3.14159", float64(3.14159)},
		{"exponent", "1.5e10", float64(1.5e10)},
		{"string", `"hello"`, "hello"},
		{"empty string", `""`, ""},
		{"escaped string", `"hello\nworld"`, "hello\nworld"},
		{"unicode", `"caf\u00e9"`, "café"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v (%T), want %v (%T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestUnmarshalArray(t *testing.T) {
	input := `[1, 2, 3, "four", true, null]`
	var result []interface{}
	err := Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	expected := []interface{}{float64(1), float64(2), float64(3), "four", true, nil}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestUnmarshalObject(t *testing.T) {
	input := `{"name": "John", "age": 30, "active": true}`
	var result map[string]interface{}
	err := Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result["name"] != "John" {
		t.Errorf("name: got %v, want John", result["name"])
	}
	if result["age"] != float64(30) {
		t.Errorf("age: got %v, want 30", result["age"])
	}
	if result["active"] != true {
		t.Errorf("active: got %v, want true", result["active"])
	}
}

func TestUnmarshalNested(t *testing.T) {
	input := `{
		"person": {
			"name": "Alice",
			"scores": [95, 87, 92]
		},
		"valid": true
	}`
	var result map[string]interface{}
	err := Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	person := result["person"].(map[string]interface{})
	if person["name"] != "Alice" {
		t.Errorf("person.name: got %v, want Alice", person["name"])
	}

	scores := person["scores"].([]interface{})
	if len(scores) != 3 {
		t.Errorf("scores length: got %d, want 3", len(scores))
	}
}

func TestUnmarshalStruct(t *testing.T) {
	type Person struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Active  bool   `json:"active"`
		private string
	}

	input := `{"name": "Bob", "age": 25, "active": true}`
	var result Person
	err := Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if result.Name != "Bob" {
		t.Errorf("Name: got %v, want Bob", result.Name)
	}
	if result.Age != 25 {
		t.Errorf("Age: got %v, want 25", result.Age)
	}
	if result.Active != true {
		t.Errorf("Active: got %v, want true", result.Active)
	}
}

func TestMarshalBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, "null"},
		{"true", true, "true"},
		{"false", false, "false"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"string", "hello", `"hello"`},
		{"escaped", "hello\nworld", `"hello\nworld"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestMarshalArray(t *testing.T) {
	input := []interface{}{1, "two", true, nil}
	result, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	expected := `[1,"two",true,null]`
	if string(result) != expected {
		t.Errorf("got %s, want %s", result, expected)
	}
}

func TestMarshalMap(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
		"num":  42,
	}
	result, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	// Keys are sorted alphabetically
	expected := `{"name":"test","num":42}`
	if string(result) != expected {
		t.Errorf("got %s, want %s", result, expected)
	}
}

func TestMarshalStruct(t *testing.T) {
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Email  string `json:"email,omitempty"`
		Ignore string `json:"-"`
	}

	input := Person{Name: "Alice", Age: 30, Ignore: "secret"}
	result, err := Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	expected := `{"name":"Alice","age":30}`
	if string(result) != expected {
		t.Errorf("got %s, want %s", result, expected)
	}
}

func TestMarshalIndent(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := MarshalIndent(input, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent error: %v", err)
	}
	expected := `{
  "name": "test"
}`
	if string(result) != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{`{"name": "test"}`, true},
		{`[1, 2, 3]`, true},
		{`"hello"`, true},
		{`42`, true},
		{`true`, true},
		{`null`, true},
		{`{invalid}`, false},
		{`[1, 2,]`, false},
		{`{"key": }`, false},
	}

	for _, tt := range tests {
		result := Valid([]byte(tt.input))
		if result != tt.valid {
			t.Errorf("Valid(%s): got %v, want %v", tt.input, result, tt.valid)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	original := Person{
		Name: "Charlie",
		Age:  35,
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	// Marshal
	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Unmarshal
	var decoded Person
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("round trip failed:\noriginal: %+v\ndecoded: %+v", original, decoded)
	}
}

func TestLexerTokens(t *testing.T) {
	input := `{"key": [1, true, null]}`
	lexer := NewLexer(input)

	expected := []TokenType{
		TokenLeftBrace,
		TokenString,
		TokenColon,
		TokenLeftBracket,
		TokenNumber,
		TokenComma,
		TokenTrue,
		TokenComma,
		TokenNull,
		TokenRightBracket,
		TokenRightBrace,
		TokenEOF,
	}

	for i, exp := range expected {
		tok := lexer.NextToken()
		if tok.Type != exp {
			t.Errorf("token %d: got %v, want %v", i, tok.Type, exp)
		}
	}
}
