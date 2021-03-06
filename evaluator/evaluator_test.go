package evaluator

import (
	"fmt"
	"github.com/qiuhoude/go-interpreter/lexer"
	"github.com/qiuhoude/go-interpreter/object"
	"github.com/qiuhoude/go-interpreter/parser"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program, object.GlobalEnv())
}

func TestEvaluator(t *testing.T) {

	Convey("TestEvalIntegerExpression", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"5", 5},
			{"10", 10},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

	Convey("TestEvalBooleanExpression", t, func() {
		cases := []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"false", false},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsBooleanObject, tt.expected)
		}
	})

	Convey("TestBangOperator", t, func() {
		cases := []struct {
			input    string
			expected bool
		}{
			{"!true", false},
			{"!false", true},
			{"!5", false},
			{"!!true", true},
			{"!!false", false},
			{"!!5", true},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsBooleanObject, tt.expected)
		}
	})

	Convey("TestMinusOrPlusOperator", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"-5", -5},
			{"-(+5)", -5},
			{"+5", 5},
			{"+(-5)", -5},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

	Convey("TestEvalIntegerExpression", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"5", 5},
			{"10", 10},
			{"-5", -5},
			{"-10", -10},
			{"5 + 5 + 5 + 5 - 10", 10},
			{"2 * 2 * 2 * 2 * 2", 32},
			{"-50 + 100 + -50", 0},
			{"5 * 2 + 10", 20},
			{"5 + 2 * 10", 25},
			{"20 + 2 * -10", 0},
			{"50 / 2 * 2 + 10", 60},
			{"2 * (5 + 10)", 30},
			{"3 * 3 * 3 + 10", 37},
			{"3 * (3 * 3) + 10", 37},
			{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

	Convey("TestEvalBooleanExpression", t, func() {
		cases := []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"false", false},
			{"1 < 2", true},
			{"1 > 2", false},
			{"1 < 1", false},
			{"1 > 1", false},
			{"1 == 1", true},
			{"1 != 1", false},
			{"1 == 2", false},
			{"1 != 2", true},

			{"true == true", true},
			{"false == false", true},
			{"true == false", false},
			{"true != false", true},
			{"false != true", true},
			{"(1 < 2) == true", true},
			{"(1 < 2) == false", false},
			{"(1 > 2) == true", false},
			{"(1 > 2) == false", true},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsBooleanObject, tt.expected)
		}
	})

	Convey("TestIfElseExpressions", t, func() {
		cases := []struct {
			input    string
			expected interface{}
		}{
			{"if (true) { 10 }", 10},
			{"if (false) { 10 }", nil},
			{"if (1) { 10 }", 10},
			{"if (1 < 2) { 10 }", 10},
			{"if (1 > 2) { 10 }", nil},
			{"if (1 > 2) { 10 } else { 20 }", 20},
			{"if (1 < 2) { 10 } else { 20 }", 10},
		}
		for _, tt := range cases {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				So(evaluated, shouldIsIntegerObject, int64(integer))
			} else {
				So(evaluated, shouldIsNullObject)
			}
		}
	})

	Convey("TestReturnStatements", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"return 10;", 10},
			{"return 10; 9;", 10},
			{"return 2 * 5; 9;", 10},
			{"9; return 2 * 5; 9;", 10},
			{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}
`, 10},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

	Convey("TestErrorHandling", t, func() {
		cases := []struct {
			input           string
			expectedMessage string
		}{
			{
				"5 + true;",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"5 + true; 5;",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"-true",
				"unknown operator: -BOOLEAN",
			},
			{
				"true + false;",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"5; true + false; 5",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"if (10 > 1) { true + false; }",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				`
if (10 > 1) {
if (10 > 1) {
return true + false;
}
return 1;
}`,
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"foobar",
				"identifier not found: foobar",
			},
			{
				`"Hello" - "World"`,
				"unknown operator: STRING - STRING",
			},
			{
				`hash{"name": "Monkey"}[fn(x) { x }];`,
				"unusable as hash key: FUNCTION",
			},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			Convey(tt.input, func() {
				So(actual, shouldIsErrorObjectMsgEq, tt.expectedMessage)
			})

		}
	})

	Convey("TestLetStatements", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"let a = 5; a;", 5},
			{"let a = 5 * 5; a;", 25},
			{"let a = 5; let b = a; b;", 5},
			{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
			{"let a = 6; if( true ){ let a = 5; }  a;", 6},
			{"let a = 10; { let a = 5; };  a;", 10},
			{"let a = 10; { a = 5; };  a;", 5},
			{"{ a = 5; };  a;", 5},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

	Convey("TestFunctionObject", t, func() {
		input := "fn(x) { x + 2; };"
		actual := testEval(input)
		So(actual, shouldIsFunctionObject)
		fn := actual.(*object.Function)
		So(len(fn.Parameters), ShouldEqual, 1)
		So(fn.Parameters[0].String(), ShouldEqual, "x")
		So(fn.Body.String(), ShouldEqual, "(x + 2)")

	})

	Convey("TestFunctionApplication", t, func() {
		cases := []struct {
			input    string
			expected int64
		}{
			{"let identity = fn(x) { x; }; identity(5);", 5},
			{"let identity = fn(x) { return x; }; identity(5);", 5},
			{"let double = fn(x) { x * 2; }; double(5);", 10},
			{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
			{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
			{"fn(x) { x; }(5)", 5},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsIntegerObject, tt.expected)
		}
	})

}
func TestBuiltinFunctions(t *testing.T) {
	Convey("TestBuiltinFunctions", t, func() {
		cases := []struct {
			input    string
			expected interface{}
		}{
			{`len("")`, 0},
			{`len("four")`, 4},
			{`len("hello world")`, 11},
			{`len(1)`, "argument to `len` not supported, got INTEGER"},
			{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				So(actual, shouldIsIntegerObject, int64(expected))
			case string:
				So(actual, shouldIsErrorObjectMsgEq, expected)
			}
		}
	})
}

func TestStringLiteral(t *testing.T) {

	Convey("TestStringLiteral", t, func() {
		cases := []struct {
			input    string
			expected string
		}{
			{`"Hello World!"`, "Hello World!"},
			{`"H@@@__@`, "H@@@__@"},
			{`"Hello" + " " + "World!"`, "Hello World!"},
			{`"Hello" +" " +  1`, "Hello 1"},
			{`"Hello" +" " + true`, "Hello true"},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			So(actual, shouldIsStringObject, tt.expected)
		}
	})
}

func TestArray(t *testing.T) {
	Convey("TestArrayLiterals", t, func() {
		input := "[1, 2 * 2, 3 + 3]"
		array := testEval(input)
		So(array, shouldIsArrayObject, 3)
		arrObj := array.(*object.Array)
		So(arrObj.Elements[0], shouldIsIntegerObject, int64(1))
		So(arrObj.Elements[1], shouldIsIntegerObject, int64(4))
		So(arrObj.Elements[2], shouldIsIntegerObject, int64(6))
	})

	Convey("TestArrayIndexExpressions", t, func() {
		cases := []struct {
			input    string
			expected interface{}
		}{
			{"[1, 2, 3][0]", 1},
			{"[1, 2, 3][1]", 2},
			{"[1, 2, 3][2]", 3},
			{"let i = 0; [1][i];", 1},
			{"[1, 2, 3][1 + 1];", 3},
			{"let myArray = [1, 2, 3]; myArray[2];", 3},
			{
				"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
				6,
			},
			{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
			{"[1, 2, 3][3]", nil},
			{"[1, 2, 3][-1]", nil},
		}
		for _, tt := range cases {
			actual := testEval(tt.input)
			switch expected := tt.expected.(type) {
			case int:
				So(actual, shouldIsIntegerObject, int64(expected))
			case nil:
				So(actual, shouldIsNullObject)
			}
		}
	})
}

func TestScript(t *testing.T) {

	Convey("TestScript", t, func() {
		cases := []struct {
			script          string
			inspectExpected interface{}
		}{
			{
				`
let map = fn(arr, f) {
      let iter = fn(arr, accumulated) {
          if (len(arr) == 0) {
              accumulated
          } else {
              iter(rest(arr), push(accumulated, f(first(arr))));
          }
      };
      iter(arr, []);
  };
let a = [1, 2, 3, 4];
let double = fn(x) { x * 2 };
map(a, double);`,
				`[2, 4, 6, 8]`,
			},
			{
				`
let reduce = fn(arr, initial, f) {
	let iter = fn(arr, result) {
		if (len(arr) == 0) {
			result
		} else {
			iter(rest(arr), f(result, first(arr)));
		}
	};
	iter(arr, initial);
};
let sum = fn(arr) {
	reduce(arr, 0, fn(initial, el) { initial + el});
};
sum([1, 2, 3, 4, 5])
`,
				`15`,
			},
		}
		for _, tt := range cases {
			eval := testEval(tt.script)
			So(eval.Inspect(), ShouldEqual, tt.inspectExpected)
		}
	})

}

func TestHashLiterals(t *testing.T) {
	Convey("TestHashLiterals", t, func() {
		input := `
let two = "two";
hash{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
}`
		evaluated := testEval(input)

		So(evaluated, shouldIsHashObjectType)
		result := evaluated.(*object.Hash)
		expected := map[object.HashKey]int64{
			(&object.String{Value: "one"}).HashKey():   1,
			(&object.String{Value: "two"}).HashKey():   2,
			(&object.String{Value: "three"}).HashKey(): 3,
			(&object.Integer{Value: 4}).HashKey():      4,
			TRUE.HashKey():                             5,
			FALSE.HashKey():                            6,
		}
		So(len(result.Pairs), ShouldEqual, len(expected))

		for expectedKey, expectedValue := range expected {
			pair, ok := result.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}
			So(pair.Value, shouldIsIntegerObject, expectedValue)
		}
	})

}

func TestHashIndexExpressions(t *testing.T) {
	Convey("TestHashIndexExpressions", t, func() {
		cases := []struct {
			input    string
			expected interface{}
		}{
			{
				`hash{"foo": 5}["foo"]`,
				5,
			},
			{
				`hash{"foo": 5}["bar"]`,
				nil,
			},
			{
				`let key = "foo"; hash{"foo": 5}[key]`,
				5,
			},
			{
				`hash{}["foo"]`,
				nil,
			},
			{
				`hash{5: 5}[5]`,
				5,
			},
			{
				`hash{true: 5}[true]`,
				5,
			},
			{
				`hash{false: 5}[false]`,
				5,
			},
		}
		for _, tt := range cases {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				So(evaluated, shouldIsIntegerObject, int64(integer))
			} else {
				So(evaluated, shouldIsNullObject)
			}
		}
	})
}

func shouldIsHashObjectType(actual interface{}, _ ...interface{}) string {
	_, ok := actual.(*object.Hash)
	if !ok {
		return fmt.Sprintf("object is not Hash. got=%T (%+v)",
			actual, actual)
	}
	return ""
}
func shouldIsArrayObject(actual interface{}, expectedList ...interface{}) string {
	// ???????????????
	expectedLen := expectedList[0].(int)

	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Sprintf("object is not Array. got=%T (%+v)",
			actual, actual)
	}
	if len(result.Elements) != expectedLen {
		return fmt.Sprintf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	return ""
}

func shouldIsStringObject(actual interface{}, expectedList ...interface{}) string {
	expected := expectedList[0].(string)

	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Sprintf("object is not String. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Sprintf("String has wrong value. expected=%q, got=%q",
			expected, result.Value)
	}
	return ""
}
func shouldIsFunctionObject(actual interface{}, _ ...interface{}) string {
	_, ok := actual.(*object.Function)
	if !ok {
		return fmt.Sprintf("object is not Function. got=%T (%+v)",
			actual, actual)
	}
	return ""
}

func shouldIsErrorObjectMsgEq(actual interface{}, expectedList ...interface{}) string {
	expected := expectedList[0].(string)

	result, ok := actual.(*object.Error)
	if !ok {
		return fmt.Sprintf("no error object returned. got=%T (%+v)",
			actual, actual)
	}
	if result.Message != expected {
		return fmt.Sprintf("wrong error message. expected=%q, got=%q",
			expected, result.Message)
	}
	return ""
}

func shouldIsNullObject(actual interface{}, _ ...interface{}) string {
	if actual != NULL {
		return fmt.Sprintf("object is not NULL. got=%T (%+v)",
			actual, actual)
	}
	return ""
}

func shouldIsBooleanObject(actual interface{}, expectedList ...interface{}) string {
	expected := expectedList[0].(bool)

	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Sprintf("object is not Boolean. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Sprintf("object has wrong value. got=%v, want=%v",
			result.Value, expected)
	}
	return ""
}

func shouldIsIntegerObject(actual interface{}, expectedList ...interface{}) string {
	expected := expectedList[0].(int64) // ????????????

	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Sprintf("object is not Integer. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Sprintf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}
	return ""
}
