package expr

import (
	"errors"
	"fmt"
	"reflect"
)

type Callable interface {
	Call(arguments []interface{}) interface{}
	ArgNum() int
}

type PrintFunc struct{}

func (PrintFunc) Call(_ *Interpreter, args []interface{}) interface{} {
	fmt.Printf("%+v", args)
	return nil
}

func (PrintFunc) ArgNum() int { return 1 }

type GoFunc struct {
	argNum int
	call   func(arguments []interface{}) interface{}
}

func (f *GoFunc) Call(args []interface{}) interface{} {
	return f.call(args)
}

func (f *GoFunc) ArgNum() int {
	return f.argNum
}

func (e *Environment) DefineGoFunc(name string, f interface{}) error {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return errors.New("not a function")
	}

	v := reflect.ValueOf(f)
	numIn := t.NumIn()
	numOut := t.NumOut()

	if numOut > 1 {
		return errors.New("too many return values")
	}

	gf := &GoFunc{
		argNum: numIn,
		call: func(arguments []interface{}) interface{} {
			in := make([]reflect.Value, 0, numIn)
			for i, inputArg := range arguments {
				argDef := t.In(i)
				inputVal := reflect.ValueOf(inputArg)

				converted, ok := TryConvert(inputVal, argDef)
				if !ok {
					return RuntimeError{msg: fmt.Sprintf("%s argument[%d] '%+v' %T is not compatible for %+v",
						name, i, inputArg, inputArg, argDef)}
				}

				in = append(in, converted)
			}

			results := v.Call(in)
			if len(results) == 0 {
				return nil
			}
			return results[0].Interface()
		},
	}

	e.Define(name, gf)
	return nil
}

// go1.17 以后可以使用 reflect.Value.CanConvert() 判断
func TryConvert(v reflect.Value, t reflect.Type) (converted reflect.Value, ok bool) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		ok = false
	}()

	converted = v.Convert(t)
	return converted, true
}
