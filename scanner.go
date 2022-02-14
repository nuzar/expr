package expr

import (
	"strconv"
	"unicode"
)

var keywords = map[string]TokenType{
	"and":   TokenAnd,
	"nil":   TokenNil,
	"or":    TokenOr,
	"true":  TokenTrue,
	"false": TokenFalse,
}

type Scanner struct {
	src    []rune
	tokens []*Token

	start, current int
	line           int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:  []rune(src),
		line: 1,
	}
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, NewToken(TokenEOF, "", nil, s.line))
	return s.tokens, nil
}

//revive:disable:cyclomatic
func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '-':
		s.addToken(TokenMinus, nil)
	case '(':
		s.addToken(TokenLeftParen, nil)
	case ')':
		s.addToken(TokenRightParen, nil)
	case '[':
		s.addToken(TokenLeftBracket, nil)
	case ']':
		s.addToken(TokenRightBracket, nil)
	case ',':
		s.addToken(TokenComma, nil)
	case '.':
		s.addToken(TokenDot, nil)
	case '!':
		if s.match('=') {
			s.addToken(TokenBangEqual, nil)
		} else {
			s.addToken(TokenBang, nil)
		}
	case '=':
		if !s.match('=') {
			return ScanError(report(s.line, "", "unexpected character"+string(c)))
		}
		s.addToken(TokenEqualEqual, nil)
	case '<':
		if s.match('=') {
			s.addToken(TokenLessEqual, nil)
		} else {
			s.addToken(TokenLess, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenGreaterEqual, nil)
		} else {
			s.addToken(TokenGreater, nil)
		}
	// case '/':
	// 	if s.match('/') {
	// 		// A comment goes until the end of the line.
	// 		for s.peek() != '\n' && !s.isAtEnd() {
	// 			s.advance()
	// 		}
	// 	} else {
	// 		s.addToken(SLASH, nil)
	// 	}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		if err := s.string(); err != nil {
			return err
		}
	default:
		if unicode.IsDigit(c) {
			s.number()
		} else if unicode.IsLetter(c) {
			s.identifier()
		} else {
			return ScanError(report(s.line, "", "unexpected character "+string(c)))
		}
	}

	return nil
}

//revive:enable:cyclomatic

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	text := string(s.src[s.start:s.current])
	s.tokens = append(s.tokens, NewToken(t, text, literal, s.line))
}

func (s *Scanner) advance() rune {
	c := s.src[s.current]
	s.current++
	return c
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.src[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.src[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.src) {
		return '\x00'
	}
	return s.src[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		// TODO: support multi-line strings?
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return ScanError(report(s.line, "", "Unterminated string."))
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.src[s.start+1 : s.current-1]
	s.addToken(TokenString, value)

	return nil
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
	}

	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	val, _ := strconv.ParseFloat(string(s.src[s.start:s.current]), 64)
	s.addToken(TokenNumber, val)
}

func (s *Scanner) identifier() {
	for isIdentifier(s.peek()) {
		s.advance()
	}

	text := string(s.src[s.start:s.current])
	if typ, ok := keywords[text]; ok {
		s.addToken(typ, nil)
	} else {
		s.addToken(TokenIdentifier, nil)
	}
}

func isIdentifier(c rune) bool {
	if isAlphaNumeric(c) {
		return true
	}
	return c == '_'
}

func isAlphaNumeric(c rune) bool {
	return unicode.IsDigit(c) || unicode.IsLetter(c)
}

type ScanError string

func (p ScanError) Error() string {
	return string(p)
}
