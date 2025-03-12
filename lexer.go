package main

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	TokenBraceL TokenType = "{"
	TokenBraceR TokenType = "}"
	TokenParenL TokenType = "("
	TokenParenR TokenType = ")"
	TokenColon  TokenType = ":"
	TokenAt     TokenType = "@" // New token for @ symbol
	TokenString TokenType = "STRING"
	TokenIdent  TokenType = "IDENT"
	TokenEOF    TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

type Lexer struct {
	input       string
	position    int
	currentChar rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = rune(l.input[l.position])
	}
	l.position++
}

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

func LexerMain() {
	lexer := NewLexer(`query GetUser { user(id: "123") { name } }`)
	for {
		tok := lexer.NextToken()
		fmt.Printf("Token: %+v\n", tok)
		if tok.Type == TokenEOF {
			break
		}
	}
}
