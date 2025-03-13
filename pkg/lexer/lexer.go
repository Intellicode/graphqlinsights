// Package lexer provides tokenization for GraphQL queries
package lexer

import (
	"unicode"
)

// TokenType represents the type of a token in the GraphQL query
type TokenType string

// Token types for GraphQL query lexing
const (
	TokenBraceL TokenType = "{"
	TokenBraceR TokenType = "}"
	TokenParenL TokenType = "("
	TokenParenR TokenType = ")"
	TokenColon  TokenType = ":"
	TokenAt     TokenType = "@" // Token for @ symbol used in directives
	TokenString TokenType = "STRING"
	TokenIdent  TokenType = "IDENT"
	TokenEOF    TokenType = "EOF"
)

// Token represents a lexical token in the GraphQL query
type Token struct {
	Type  TokenType
	Value string
}

// Helper functions for character classification
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// Lexer represents a lexical analyzer for GraphQL queries
type Lexer struct {
	input       string
	position    int
	currentChar rune
}

// NewLexer creates a new lexer for the given input string
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input string
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = rune(l.input[l.position])
	}
	l.position++
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	for unicode.IsSpace(l.currentChar) {
		l.readChar()
	}

	switch l.currentChar {
	case '{':
		l.readChar()
		return Token{TokenBraceL, "{"}
	case '}':
		l.readChar()
		return Token{TokenBraceR, "}"}
	case '(':
		l.readChar()
		return Token{TokenParenL, "("}
	case ')':
		l.readChar()
		return Token{TokenParenR, ")"}
	case ':':
		l.readChar()
		return Token{TokenColon, ":"}
	case '@': // Handle @ symbol for directives
		l.readChar()
		return Token{TokenAt, "@"}
	case '"':
		l.readChar()
		start := l.position - 1
		for l.currentChar != '"' {
			l.readChar()
		}
		value := l.input[start:l.position]
		l.readChar()
		return Token{TokenString, value}
	case 0:
		return Token{TokenEOF, ""}
	default:
		if isLetter(l.currentChar) {
			start := l.position - 1
			for isLetter(l.currentChar) || isDigit(l.currentChar) {
				l.readChar()
			}
			return Token{TokenIdent, l.input[start : l.position-1]}
		}
	}
	return Token{TokenEOF, ""}
}
