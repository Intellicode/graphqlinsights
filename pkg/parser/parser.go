// Package parser provides parsing functionality for GraphQL queries
package parser

import (
	"fmt"
	"strings"

	"github.com/tom/graphqlinsights/pkg/lexer"
)

// NodeType represents the type of a node in the GraphQL AST
type NodeType string

// Node types for GraphQL query parsing
const (
	NodeQuery     NodeType = "Query"
	NodeField     NodeType = "Field"
	NodeDirective NodeType = "Directive" // Node type for directives
)

// Node represents a node in the GraphQL AST
type Node struct {
	Type         NodeType
	Name         string
	Arguments    map[string]string
	Directives   []*Node // Field for directives
	SelectionSet []*Node
}

// Print returns a string representation of the node with proper indentation
func (n *Node) Print(indent string) string {
	result := fmt.Sprintf("%s%s: %s\n", indent, n.Type, n.Name)

	for argName, argValue := range n.Arguments {
		result += fmt.Sprintf("%s  Arg: %s = %s\n", indent, argName, argValue)
	}

	for _, directive := range n.Directives {
		result += fmt.Sprintf("%s  Directive: @%s\n", indent, directive.Name)
		for argName, argValue := range directive.Arguments {
			result += fmt.Sprintf("%s    Arg: %s = %s\n", indent, argName, argValue)
		}
	}

	for _, child := range n.SelectionSet {
		result += child.Print(indent + "  ")
	}

	return result
}

// Parser represents a parser for GraphQL queries
type Parser struct {
	lexer *lexer.Lexer
	curr  lexer.Token
}

// NewParser creates a new parser for the given input string
func NewParser(input string) *Parser {
	lex := lexer.NewLexer(input)
	return &Parser{lexer: lex, curr: lex.NextToken()}
}

// eat consumes the current token if it matches the expected type
func (p *Parser) eat(t lexer.TokenType) {
	if p.curr.Type == t {
		p.curr = p.lexer.NextToken()
	} else {
		panic(fmt.Sprintf("Unexpected token: expected %s but got %s", t, p.curr.Type))
	}
}

// ParseDirective parses a directive in a GraphQL query
func (p *Parser) ParseDirective() *Node {
	p.eat(lexer.TokenAt)
	name := p.curr.Value
	p.eat(lexer.TokenIdent)

	args := make(map[string]string)
	if p.curr.Type == lexer.TokenParenL {
		p.eat(lexer.TokenParenL)
		// Parse one or more arguments
		for p.curr.Type == lexer.TokenIdent {
			argName := p.curr.Value
			p.eat(lexer.TokenIdent)
			p.eat(lexer.TokenColon)
			argValue := p.curr.Value
			p.eat(lexer.TokenString)
			// Strip quotes from string values
			argValue = strings.Trim(argValue, "\"")
			args[argName] = argValue
			// If there are more arguments, they need to be separated properly
			// In a more complete implementation, we would handle commas here
		}
		p.eat(lexer.TokenParenR)
	}

	return &Node{
		Type:      NodeDirective,
		Name:      name,
		Arguments: args,
	}
}

// ParseField parses a field in a GraphQL query
func (p *Parser) ParseField() *Node {
	name := p.curr.Value
	p.eat(lexer.TokenIdent)

	args := make(map[string]string)
	if p.curr.Type == lexer.TokenParenL {
		p.eat(lexer.TokenParenL)
		argName := p.curr.Value
		p.eat(lexer.TokenIdent)
		p.eat(lexer.TokenColon)
		argValue := p.curr.Value
		p.eat(lexer.TokenString)
		// Strip quotes from string values
		argValue = strings.Trim(argValue, "\"")
		args[argName] = argValue
		p.eat(lexer.TokenParenR)
	}

	// Parse directives if present
	var directives []*Node
	for p.curr.Type == lexer.TokenAt {
		directives = append(directives, p.ParseDirective())
	}

	var selectionSet []*Node
	if p.curr.Type == lexer.TokenBraceL {
		p.eat(lexer.TokenBraceL)
		for p.curr.Type == lexer.TokenIdent {
			selectionSet = append(selectionSet, p.ParseField())
		}
		p.eat(lexer.TokenBraceR)
	}

	return &Node{
		Type:         NodeField,
		Name:         name,
		Arguments:    args,
		Directives:   directives,
		SelectionSet: selectionSet,
	}
}

// ParseQuery parses a GraphQL query
func (p *Parser) ParseQuery() *Node {
	p.eat(lexer.TokenIdent) // eat "query"
	name := p.curr.Value
	p.eat(lexer.TokenIdent)

	// Parse directives at query level if present
	var directives []*Node
	for p.curr.Type == lexer.TokenAt {
		directives = append(directives, p.ParseDirective())
	}

	p.eat(lexer.TokenBraceL)
	var selectionSet []*Node
	for p.curr.Type == lexer.TokenIdent {
		selectionSet = append(selectionSet, p.ParseField())
	}
	p.eat(lexer.TokenBraceR)

	return &Node{
		Type:         NodeQuery,
		Name:         name,
		Directives:   directives,
		SelectionSet: selectionSet,
	}
}
