package expr

import (
	"fmt"
)

type TokenType int

const (
	// Single-character tokens.
	TokenLeftParen    TokenType = iota // (
	TokenRightParen                    // )
	TokenLeftBracket                   // [
	TokenRightBracket                  // ]
	TokenComma                         // ,
	TokenDot                           // .

	// One or two character tokens.
	TokenMinus        // -
	TokenBang         // !
	TokenBangEqual    // !=
	TokenEqualEqual   // ==
	TokenGreater      // >
	TokenGreaterEqual // >=
	TokenLess         // <
	TokenLessEqual    // <=

	// Literals.
	TokenIdentifier // a
	TokenString     // "123"
	TokenNumber     // 123

	// Keywords.
	TokenAnd   // and
	TokenOr    // or
	TokenNil   // nil
	TokenTrue  // true
	TokenFalse // false

	TokenEOF
)

type Token struct {
	typ     TokenType
	lexeme  string
	literal interface{}
	line    int
}

func NewToken(typ TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		typ:     typ,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func (t Token) string() string {
	return fmt.Sprintf("%d %s %v", t.typ, t.lexeme, t.literal)
}
