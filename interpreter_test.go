package expr

import (
	"testing"
)

func toExpr(src string) (Expr, error) {
	s := NewScanner(src)
	tokens, err := s.ScanTokens()
	if err != nil {
		return nil, err
	}

	p := NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func Test_equal(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `1 == 1`, expect: true},
		{src: `1 != 1`, expect: false},
		{src: `1 != 2`, expect: true},
		{src: `"a" == "a"`, expect: true},
		{src: `"a" != "a"`, expect: false},
		{src: `"a" != "b"`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_call(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `ner_entities("product_type") == "面膜"`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_comparison(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"yes":          true,
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `ner_entities("age") == 18`, expect: true},
		{src: `ner_entities("age") >= 18`, expect: true},
		{src: `ner_entities("age") > 17`, expect: true},
		{src: `ner_entities("age") <= 18`, expect: true},
		{src: `ner_entities("age") < 19`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_array(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"yes":          true,
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `ner_entities("efficacy") == ["补水", "抗皱"]`, expect: true},
		{src: `ner_entities("efficacy") == ["补水", "+皱"]`, expect: false},
		{src: `ner_entities("efficacy") != ["补水", "+皱"]`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_bool(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"yes":          true,
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `ner_entities("yes")`, expect: true},
		{src: `ner_entities("yes") == true`, expect: true},
		{src: `ner_entities("yes") != false`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_grouping(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"yes":          true,
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `true or (false or true)`, expect: true},
		{src: `false or (false and true)`, expect: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_unary(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"yes":          true,
	}
	closure := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", closure); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name   string
		src    string
		expect bool
	}{
		{src: `!true`, expect: false},
		{src: `!false`, expect: true},
		{src: `-1 < 0`, expect: true},
		{src: `1 > 0`, expect: true},
		{src: `--1 == 1`, expect: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}
		})
	}
}

func Test_logical(t *testing.T) {
	p := NewInterpreter()

	data := map[string]interface{}{
		"product_type": "面膜",
		"skin_color":   "黑色",
		"age":          18,
		"efficacy":     []string{"补水", "抗皱"},
		"ruok":         true,
	}
	nerEntities := func(name string) interface{} {
		return data[name]
	}
	if err := p.Environment.DefineGoFunc("ner_entities", nerEntities); err != nil {
		t.Fatal(err)
	}

	cnt := 0
	callCounter := func() bool {
		cnt++
		return true
	}
	if err := p.Environment.DefineGoFunc("call_counter", callCounter); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name          string
		src           string
		expect        bool
		checkCounter  bool
		expectCounter int
	}{
		{src: `ner_entities("ruok")`, expect: true},
		{src: `!ner_entities("ruok")`, expect: false},
		{src: `ner_entities("ruok") == true`, expect: true},
		{src: `ner_entities("ruok") != false`, expect: true},
		{src: `true or call_counter()`, expect: true, checkCounter: true, expectCounter: 0},
		{src: `false or call_counter()`, expect: true, checkCounter: true, expectCounter: 1},
		{src: `true and call_counter()`, expect: true, checkCounter: true, expectCounter: 1},
		{src: `false and call_counter()`, expect: false, checkCounter: true, expectCounter: 0},
		{src: `true or false`, expect: true},
		{src: `true and false`, expect: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				cnt = 0
			}()

			e, err := toExpr(tc.src)
			if err != nil {
				t.Logf("parse expr failed: %s", err)
				t.FailNow()
			}

			res, err := p.Interpret(e)
			if err != nil {
				t.Logf("interpret expr failed: %s", err)
				t.FailNow()
			}

			b, ok := res.(bool)
			if !ok {
				t.Logf("result is not bool: %T %+v", res, res)
				t.FailNow()
			}
			if b != tc.expect {
				t.Logf("expect %v, got %v", tc.expect, b)
				t.FailNow()
			}

			if tc.checkCounter {
				if cnt != tc.expectCounter {
					t.Logf("expect counter %d, got %d", tc.expectCounter, cnt)
					t.FailNow()
				}
			}
		})
	}
}

func Test_runtime_error(t *testing.T) {
	p := NewInterpreter()

	die := func() interface{} {
		panic("die")
	}

	if err := p.Environment.DefineGoFunc("die", die); err != nil {
		t.Fatal(err)
	}

	src := `die()`

	e, err := toExpr(src)
	if err != nil {
		t.Logf("parse expr failed: %s", err)
		t.FailNow()
	}

	res, err := p.Interpret(e)
	_, ok := err.(RuntimeError)
	if !ok {
		t.Fatal("want RuntimeError")
	}

	if res != nil {
		t.Fatalf("want nil result, got %+v", res)
	}
}
