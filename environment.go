package expr

import (
	"fmt"
)

// Environment stores global symbols

type Environment struct {
	values map[string]interface{}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name *Token) (interface{}, error) {
	v, ok := e.values[name.lexeme]
	if ok {
		return v, nil
	}

	return nil, RuntimeError{msg: fmt.Sprintf("undefined symbol %s", name.lexeme)}
}
