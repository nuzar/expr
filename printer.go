package expr

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

var _ ExprVisitorStr = (*AstPrinter)(nil)

func (p *AstPrinter) Print(expr Expr) string {
	return expr.AcceptStr(p)
}

func (p *AstPrinter) VisitExprBinaryStr(expr *ExprBinary) string {
	return p.block(expr.operator.lexeme, "(", ")", expr.left, expr.right)
}

func (p *AstPrinter) VisitExprGroupingStr(expr *ExprGrouping) string {
	return p.block("group", "(", ")", expr.expression)
}

func (p *AstPrinter) VisitExprLiteralStr(expr *ExprLiteral) string {
	if expr.value == nil {
		return "nil"
	}

	return fmt.Sprint(expr.value)
}

func (p *AstPrinter) VisitExprUnaryStr(expr *ExprUnary) string {
	return p.block(expr.operator.lexeme, "(", ")", expr.right)
}

func (p *AstPrinter) VisitExprCallStr(expr *ExprCall) string {
	return p.block(expr.callee.AcceptStr(p), "(", ")", expr.arguments...)
}

func (p *AstPrinter) VisitExprLogicalStr(expr *ExprLogical) string {
	return p.block(expr.operator.lexeme, "(", ")", expr.left, expr.right)
}

func (p *AstPrinter) VisitExprVariableStr(expr *ExprVariable) string {
	return expr.name.lexeme
}

func (p *AstPrinter) VisitExprArrayStr(expr *ExprArray) string {
	return p.block(expr.bracket.lexeme, "[", "]", expr.items...)
}

func (p *AstPrinter) block(name, start, end string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString(start)
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.AcceptStr(p))
	}
	builder.WriteString(end)

	return builder.String()
}
