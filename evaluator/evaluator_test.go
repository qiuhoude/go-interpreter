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
	expected := expectedList[0].(int64) // 预期结果

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
