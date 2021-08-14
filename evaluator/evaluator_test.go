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
	return Eval(program)
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
}

func shouldIsBooleanObject(actual interface{}, expectedList ...interface{}) string {
	expected := expectedList[0].(bool)

	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Sprintf("object is not Integer. got=%T (%+v)",
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
		return fmt.Sprintf("object is not Boolean. got=%T (%+v)",
			actual, actual)
	}
	if result.Value != expected {
		return fmt.Sprintf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}
	return ""
}
