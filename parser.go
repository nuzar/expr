package expr

import (
	"fmt"
	"reflect"
)

/*
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
*/

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (expr Expr, err error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.or()
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(TokenOr) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = NewExprLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(TokenAnd) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = NewExprLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(TokenBangEqual, TokenEqualEqual) {
		var operator = p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewExprBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	// 按完整的实现，下一级应该是term
	nextLevel := p.unary

	expr, err := nextLevel()
	if err != nil {
		return nil, err
	}

	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		var operator = p.previous()
		right, err := nextLevel()
		if err != nil {
			return nil, err
		}
		expr = NewExprBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match() { // MINUS, PLUS
		var operator = p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewExprBinary(expr, operator, right)
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match() { // SLASH, STAR
		var operator = p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewExprBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(TokenBang, TokenMinus) {
		var operator = p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewExprUnary(operator, right), nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.match(TokenLeftParen) {
		expr, err = p.finishCall(expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr
	if !p.check(TokenRightParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, e)
		for p.match(TokenComma) {
			// TODO: limit arguments size
			e, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, e)
		}
	}
	paren, err := p.consume(TokenRightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return NewExprCall(callee, paren, arguments), nil
}

func (p *Parser) finishArray() (Expr, error) {
	var items []Expr
	if !p.check(TokenRightBracket) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		items = append(items, e)
		for p.match(TokenComma) {
			e, err := p.expression()
			if err != nil {
				return nil, err
			}
			items = append(items, e)
		}
	}
	bracket, err := p.consume(TokenRightBracket, "expect ']'")
	if err != nil {
		return nil, err
	}

	return NewExprArray(bracket, items), nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TokenFalse) {
		return NewExprLiteral(false, reflect.Bool), nil
	}
	if p.match(TokenTrue) {
		return NewExprLiteral(true, reflect.Bool), nil
	}
	if p.match(TokenNumber) {
		return NewExprLiteral(p.previous().literal, reflect.Float64), nil
	}
	if p.match(TokenString) {
		return NewExprLiteral(p.previous().literal, reflect.String), nil
	}
	if p.match(TokenIdentifier) {
		return NewExprVariable(p.previous()), nil
	}
	if p.match(TokenLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(TokenRightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewExprGrouping(expr), nil
	}

	if p.match(TokenLeftBracket) {
		expr, err := p.finishArray()
		if err != nil {
			return nil, err
		}
		return expr, nil
	}

	return nil, p.Error(p.peek(), "Expect expression.")
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().typ == t
}

func (p *Parser) isAtEnd() bool {
	return p.peek().typ == TokenEOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
func (p *Parser) consume(typ TokenType, msg string) (*Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	return nil, p.Error(p.peek(), msg)
}

func (p *Parser) Error(token *Token, msg string) ParseError {
	var s string
	if token.typ == TokenEOF {
		s = report(token.line, " at end", msg)
	} else {
		s = report(token.line, " at '"+token.lexeme+"'", msg)
	}
	return ParseError(s)
}

func report(line int, where string, msg string) string {
	return fmt.Sprintf("[line %d ] %s: %s\n", line, where, msg)
}

type ParseError string

func (p ParseError) Error() string {
	return string(p)
}
