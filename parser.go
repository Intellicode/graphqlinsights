package main

import (
	"fmt"
	"strings"
)

type NodeType string

const (
	NodeQuery NodeType = "Query"
	NodeField NodeType = "Field"
)

type Node struct {
	Type         NodeType
	Name         string
	Arguments    map[string]string
	SelectionSet []*Node
}

func (n *Node) Print(indent string) string {
	result := fmt.Sprintf("%s%s: %s\n", indent, n.Type, n.Name)
	for argName, argValue := range n.Arguments {
		result += fmt.Sprintf("%s  Arg: %s = %s\n", indent, argName, argValue)
	}
	for _, child := range n.SelectionSet {
		result += child.Print(indent + "  ")
	}
	return result
}

type Parser struct {
	lexer *Lexer
	curr  Token
}

func NewParser(input string) *Parser {
	lexer := NewLexer(input)
	return &Parser{lexer: lexer, curr: lexer.NextToken()}
}

func (p *Parser) eat(t TokenType) {
	if p.curr.Type == t {
		p.curr = p.lexer.NextToken()
	} else {
		panic(fmt.Sprintf("Unexpected token: expected %s but got %s", t, p.curr.Type))
	}
}

func (p *Parser) parseField() *Node {
	name := p.curr.Value
	p.eat(TokenIdent)

	args := make(map[string]string)
	if p.curr.Type == TokenParenL {
		p.eat(TokenParenL)
		argName := p.curr.Value
		p.eat(TokenIdent)
		p.eat(TokenColon)
		argValue := p.curr.Value
		p.eat(TokenString)
		// Strip quotes from string values
		argValue = strings.Trim(argValue, "\"")
		args[argName] = argValue
		p.eat(TokenParenR)
	}

	var selectionSet []*Node
	if p.curr.Type == TokenBraceL {
		p.eat(TokenBraceL)
		for p.curr.Type == TokenIdent {
			selectionSet = append(selectionSet, p.parseField())
		}
		p.eat(TokenBraceR)
	}

	return &Node{Type: NodeField, Name: name, Arguments: args, SelectionSet: selectionSet}
}

func (p *Parser) parseQuery() *Node {
	p.eat(TokenIdent) // eat "query"
	name := p.curr.Value
	p.eat(TokenIdent)
	p.eat(TokenBraceL)

	var selectionSet []*Node
	for p.curr.Type == TokenIdent {
		selectionSet = append(selectionSet, p.parseField())
	}
	p.eat(TokenBraceR)

	return &Node{Type: NodeQuery, Name: name, SelectionSet: selectionSet}
}

// ParseMain demonstrates how the input of the lexer goes to the parser
func ParseMain() {
	input := `query GetUser { user(id: "123") { name } }`
	lexer := NewLexer(input)
	parser := &Parser{lexer: lexer, curr: lexer.NextToken()}
	parsedQuery := parser.parseQuery()
	fmt.Print(parsedQuery.Print(""))
}
