package main

import (
	"fmt"
	"log"
	"testing"
)

func TestParseQuery(t *testing.T) {
	log.Println("Starting TestParseQuery")
	tests := []struct {
		name  string
		input string
		want  *Node
	}{
		{
			name:  "Simple query",
			input: `query GetUser { user(id: "123") { name } }`,
			want: &Node{
				Type: NodeQuery,
				Name: "GetUser",
				SelectionSet: []*Node{
					{
						Type:      NodeField,
						Name:      "user",
						Arguments: map[string]string{"id": "123"},
						SelectionSet: []*Node{
							{Type: NodeField, Name: "name"},
						},
					},
				},
			},
		},
		{
			name:  "Nested query",
			input: `query GetUser { user(id: "123") { name friends { name } } }`,
			want: &Node{
				Type: NodeQuery,
				Name: "GetUser",
				SelectionSet: []*Node{
					{
						Type:      NodeField,
						Name:      "user",
						Arguments: map[string]string{"id": "123"},
						SelectionSet: []*Node{
							{Type: NodeField, Name: "name"},
							{
								Type: NodeField,
								Name: "friends",
								SelectionSet: []*Node{
									{Type: NodeField, Name: "name"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("Running test: %s", tt.name)
			lexer := NewLexer(tt.input)
			parser := &Parser{lexer: lexer, curr: lexer.NextToken()}
			parsedQuery := parser.parseQuery()

			// Use our custom compareNodes function to compare node structures
			if !compareNodes(parsedQuery, tt.want) {
				t.Errorf("Node structures not equal")
				t.Logf("Detailed comparison result: %s", detailedCompare(parsedQuery, tt.want))
			}
		})
	}
}

// Helper function for detailed comparison and debugging
func detailedCompare(got, want *Node) string {
	if got == nil && want == nil {
		return "Both nodes are nil"
	}
	if got == nil {
		return "Got node is nil, want is not nil"
	}
	if want == nil {
		return "Want node is nil, got is not nil"
	}

	result := ""

	// Compare basic properties
	if got.Type != want.Type {
		result += fmt.Sprintf("Type mismatch: got %s, want %s\n", got.Type, want.Type)
	}
	if got.Name != want.Name {
		result += fmt.Sprintf("Name mismatch: got %s, want %s\n", got.Name, want.Name)
	}

	// Compare arguments
	if len(got.Arguments) != len(want.Arguments) {
		result += fmt.Sprintf("Arguments length mismatch: got %d, want %d\n", len(got.Arguments), len(want.Arguments))
	} else {
		for k, v := range got.Arguments {
			if wantVal, ok := want.Arguments[k]; !ok {
				result += fmt.Sprintf("Missing argument in want: %s\n", k)
			} else if wantVal != v {
				result += fmt.Sprintf("Argument value mismatch for %s: got %s, want %s\n", k, v, wantVal)
			}
		}
		for k := range want.Arguments {
			if _, ok := got.Arguments[k]; !ok {
				result += fmt.Sprintf("Missing argument in got: %s\n", k)
			}
		}
	}

	// Compare selection sets
	if len(got.SelectionSet) != len(want.SelectionSet) {
		result += fmt.Sprintf("SelectionSet length mismatch: got %d, want %d\n", len(got.SelectionSet), len(want.SelectionSet))
	} else {
		for i := range got.SelectionSet {
			childResult := detailedCompare(got.SelectionSet[i], want.SelectionSet[i])
			if childResult != "" {
				result += fmt.Sprintf("SelectionSet[%d] differences:\n%s", i, childResult)
			}
		}
	}

	if result == "" {
		return "Nodes are structurally equivalent but have different memory addresses"
	}
	return result
}

// Function to compare node structures properly
func compareNodes(got, want *Node) bool {
	// Check for nil nodes
	if got == nil && want == nil {
		return true
	}
	if got == nil || want == nil {
		return false
	}

	// Compare basic properties
	if got.Type != want.Type || got.Name != want.Name {
		return false
	}

	// Compare arguments
	if len(got.Arguments) != len(want.Arguments) {
		return false
	}
	for k, v := range got.Arguments {
		wantVal, ok := want.Arguments[k]
		if !ok || wantVal != v {
			return false
		}
	}

	// Compare selection sets
	if len(got.SelectionSet) != len(want.SelectionSet) {
		return false
	}
	for i := range got.SelectionSet {
		if !compareNodes(got.SelectionSet[i], want.SelectionSet[i]) {
			return false
		}
	}
	return true
}
