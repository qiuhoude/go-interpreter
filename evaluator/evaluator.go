package evaluator

import (
	"github.com/qiuhoude/go-interpreter/ast"
	"github.com/qiuhoude/go-interpreter/object"
	"github.com/qiuhoude/go-interpreter/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statements
	case *ast.Program: //AST root node
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
		// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	}
	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "+":
		return evalMinusOrPlusOperatorExpression(token.PLUS, right)
	case "-":
		return evalMinusOrPlusOperatorExpression(token.MINUS, right)
	default:
		return NULL

	}
}

func evalMinusOrPlusOperatorExpression(op token.TokenType, right object.Object) object.Object {
	r, ok := right.(*object.Integer)
	if !ok {
		return NULL
	}
	if op == token.MINUS {
		r.Value = -r.Value
	}
	return r
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, s := range stmts { // 解析最后一条语句才是返回值
		result = Eval(s)
	}

	return result
}

/*
pseudocode
function eval(astNode) {
	if (astNode is integerliteral) {
		return astNode.integerValue
	} else if (astNode is booleanLiteral) {
		return astNode.booleanValue
	} else if (astNode is infixExpression) {
		leftEvaluated = eval(astNode.Left)
		rightEvaluated = eval(astNode.Right)
	if astNode.Operator == "+" {
		return leftEvaluated + rightEvaluated
	} else if ast.Operator == "-" {
		return leftEvaluated - rightEvaluated
	}
	}
}
*/
