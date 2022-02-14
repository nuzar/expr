package expr

import (
	"fmt"
	"reflect"
	"runtime/debug"
)

func Run(src string) (interface{}, error) {
	var scanner = NewScanner(src)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	p := NewInterpreter()
	return p.Interpret(expr)
}

type Interpreter struct {
	Environment *Environment
}

var _ ExprVisitorObj = (*Interpreter)(nil)

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{
		Environment: &Environment{
			values: make(map[string]interface{}),
		},
	}

	return interpreter
}

func (p *Interpreter) Interpret(expr Expr) (res interface{}, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		switch tr := r.(type) {
		case error:
			err = tr
		default:
			stackInfo := debug.Stack()
			err = RuntimeError{msg: fmt.Sprintf("runtime err: %s, stack info: %s", tr, string(stackInfo))}
		}
	}()

	return p.evaluate(expr)
}

func (p *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.AcceptObj(p)
}

const s string = ""
const f float64 = 0

var strType = reflect.TypeOf(s)
var numType = reflect.TypeOf(f)

func (p *Interpreter) VisitExprLiteralObj(expr *ExprLiteral) (interface{}, error) {
	rv := reflect.ValueOf(expr.value)

	switch expr.rtype {
	case reflect.Bool:
		return expr.value, nil
	case reflect.String:
		return rv.Convert(strType).Interface(), nil
	case reflect.Float64:
		return rv.Convert(numType).Interface(), nil
	}

	return expr.value, nil
}

func (p *Interpreter) VisitExprLogicalObj(expr *ExprLogical) (interface{}, error) {
	left, err := p.evaluate(expr.left)
	if err != nil {
		return false, err
	}

	bLeft, err := isTruthy(left)
	if err != nil {
		return false, err
	}

	if expr.operator.typ == TokenOr {
		if bLeft {
			return true, nil
		}
	} else {
		if !bLeft {
			return false, nil
		}
	}

	right, err := p.evaluate(expr.right)
	if err != nil {
		return false, err
	}

	return isTruthy(right)
}

func (p *Interpreter) VisitExprGroupingObj(expr *ExprGrouping) (interface{}, error) {
	return p.evaluate(expr.expression)
}

func (p *Interpreter) VisitExprUnaryObj(expr *ExprUnary) (interface{}, error) {
	right, err := p.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.typ {
	case TokenMinus:
		r, isNumber := toNumber(right)
		if !isNumber {
			return nil, RuntimeError{msg: fmt.Sprintf("%+v (%T) is not number", right, right)}
		}
		return -r, nil
	case TokenBang:
		res, err := isTruthy(right)
		if err != nil {
			return nil, err
		}
		return !res, nil
	default:
		return nil, RuntimeError{msg: fmt.Sprintf("unknown operator %s", expr.operator.lexeme)}
	}
}

func (p *Interpreter) VisitExprVariableObj(expr *ExprVariable) (interface{}, error) {
	return p.lookup(expr.name)
}

func (p *Interpreter) lookup(name *Token) (interface{}, error) {
	return p.Environment.Get(name)
}

//revive:disable:cyclomatic
func (p *Interpreter) VisitExprBinaryObj(expr *ExprBinary) (interface{}, error) {
	left, err := p.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	right, err := p.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	ln, lIsNumber := toNumber(left)
	rn, rIsNumber := toNumber(right)

	if lIsNumber != rIsNumber {
		return nil, RuntimeError{msg: fmt.Sprintf("%+v %s %+v is not number",
			left, expr.operator.lexeme, right)}
	}

	switch expr.operator.typ {
	// +-*/
	case TokenGreater:
		return ln > rn, nil
	case TokenGreaterEqual:
		return ln >= rn, nil
	case TokenLess:
		return ln < rn, nil
	case TokenLessEqual:
		return ln <= rn, nil
	case TokenEqualEqual:
		if lIsNumber {
			return p.isEqual(ln, rn)
		}
		return p.isEqual(left, right)
	case TokenBangEqual:
		var eq bool
		var err error
		if lIsNumber {
			eq, err = p.isEqual(ln, rn)
		} else {
			eq, err = p.isEqual(left, right)
		}
		return !eq, err
	default:
		return nil, RuntimeError{msg: fmt.Sprintf("unknown operator %s", expr.operator.lexeme)}
	}
}

//revive:enable:cyclomatic

func (p *Interpreter) VisitExprCallObj(expr *ExprCall) (interface{}, error) {
	callee, err := p.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}

	callable, ok := callee.(Callable)
	if !ok {
		return nil, RuntimeErrWithToken(expr.paren, "not callable")
	}

	if len(expr.arguments) != callable.ArgNum() {
		return nil, RuntimeErrWithToken(expr.paren,
			fmt.Sprintf("want %d but got %d arguments", callable.ArgNum(), len(expr.arguments)))
	}

	arguments := make([]interface{}, 0, len(expr.arguments))
	for _, arg := range expr.arguments {
		argV, err := p.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argV)
	}

	return callable.Call(arguments), nil
}

func (p *Interpreter) VisitExprArrayObj(expr *ExprArray) (interface{}, error) {
	items := make([]interface{}, 0, len(expr.items))
	for _, item := range expr.items {
		items = append(items, item)
	}
	return items, nil
}

func isTruthy(obj interface{}) (bool, error) {
	if obj == nil {
		return false, RuntimeError{"nil value"}
	}

	b, ok := obj.(bool)
	if !ok {
		return false, RuntimeError{fmt.Sprintf("not bool value: %+v", obj)}
	}

	return b, nil
}

// https://golang.google.cn/ref/spec#Comparison_operators
func (p *Interpreter) isEqual(a, b interface{}) (res bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = false
		}
	}()

	isSlice := reflect.TypeOf(a).Kind() == reflect.Slice
	if isSlice {
		return p.isSliceEqual(a, b)
	}

	return a == b, nil
}

func (p *Interpreter) isSliceEqual(a, b interface{}) (res bool, err error) {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Len() != vb.Len() {
		return false, nil
	}

	length := va.Len()
	for i := 0; i < length; i++ {
		itema := va.Index(i)
		itemb := vb.Index(i)

		converted, ok := TryConvert(itema, itemb.Type())
		if !ok {
			return false, nil
		}

		ai := converted.Interface()
		bi := itemb.Interface()

		if ae, ok := ai.(Expr); ok {
			if ai, err = p.evaluate(ae); err != nil {
				return false, err
			}
		}
		if be, ok := bi.(Expr); ok {
			if bi, err = p.evaluate(be); err != nil {
				return false, err
			}
		}

		if ai != bi {
			return false, nil
		}
	}

	return true, nil
}

func toNumber(obj interface{}) (float64, bool) {
	switch n := obj.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

type RuntimeError struct {
	msg string
}

func RuntimeErrWithToken(t *Token, msg string) RuntimeError {
	return RuntimeError{msg: fmt.Sprintf("%s: %s", t.lexeme, msg)}
}

func (e RuntimeError) Error() string {
	return e.msg
}
