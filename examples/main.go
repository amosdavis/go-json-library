package main

import (
	"fmt"
	"log"

	json "github.com/adavi068/go-json-library"
)

// User demonstrates struct tags for JSON field mapping
type User struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email,omitempty"`
	Age       int      `json:"age"`
	Active    bool     `json:"active"`
	Roles     []string `json:"roles"`
	private   string   // unexported fields are ignored
	APIKey    string   `json:"-"` // explicitly ignored
}

// Address demonstrates nested structs
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type Company struct {
	Name      string  `json:"name"`
	Address   Address `json:"address"`
	Employees int     `json:"employees"`
}

func main() {
	// Example 1: Marshal a struct to JSON
	fmt.Println("=== Example 1: Marshal Struct ===")
	user := User{
		ID:     1,
		Name:   "Alice Smith",
		Email:  "alice@example.com",
		Age:    28,
		Active: true,
		Roles:  []string{"admin", "developer"},
		APIKey: "secret-key-123", // will be ignored due to json:"-"
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	fmt.Println()

	// Example 2: Pretty print with MarshalIndent
	fmt.Println("=== Example 2: Pretty Print ===")
	prettyData, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(prettyData))
	fmt.Println()

	// Example 3: Unmarshal JSON into a struct
	fmt.Println("=== Example 3: Unmarshal to Struct ===")
	jsonStr := `{
		"id": 2,
		"name": "Bob Johnson",
		"age": 35,
		"active": false,
		"roles": ["viewer"]
	}`

	var user2 User
	if err := json.Unmarshal([]byte(jsonStr), &user2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed User: %+v\n", user2)
	fmt.Println()

	// Example 4: Work with nested structures
	fmt.Println("=== Example 4: Nested Structures ===")
	company := Company{
		Name: "Tech Corp",
		Address: Address{
			Street:  "123 Innovation Way",
			City:    "San Francisco",
			Country: "USA",
		},
		Employees: 500,
	}

	companyJSON, _ := json.MarshalIndent(company, "", "  ")
	fmt.Println(string(companyJSON))
	fmt.Println()

	// Example 5: Parse into map[string]interface{} for dynamic JSON
	fmt.Println("=== Example 5: Dynamic JSON Parsing ===")
	dynamicJSON := `{
		"type": "event",
		"data": {
			"id": 12345,
			"tags": ["important", "urgent"],
			"metadata": null
		}
	}`

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(dynamicJSON), &result); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type: %v\n", result["type"])
	if data, ok := result["data"].(map[string]interface{}); ok {
		fmt.Printf("Data ID: %v\n", data["id"])
		fmt.Printf("Tags: %v\n", data["tags"])
	}
	fmt.Println()

	// Example 6: Validate JSON
	fmt.Println("=== Example 6: JSON Validation ===")
	validJSON := `{"name": "test", "value": 123}`
	invalidJSON := `{"name": "test", value: 123}` // missing quotes around value

	fmt.Printf("Valid JSON:   %v -> %v\n", validJSON, json.Valid([]byte(validJSON)))
	fmt.Printf("Invalid JSON: %v -> %v\n", invalidJSON, json.Valid([]byte(invalidJSON)))
	fmt.Println()

	// Example 7: omitempty tag behavior
	fmt.Println("=== Example 7: omitempty Behavior ===")
	userNoEmail := User{
		ID:     3,
		Name:   "Charlie",
		Age:    22,
		Active: true,
		Roles:  []string{},
	}

	noEmailJSON, _ := json.MarshalIndent(userNoEmail, "", "  ")
	fmt.Println(string(noEmailJSON))
	fmt.Println("(Note: 'email' field is omitted because it's empty)")
}
