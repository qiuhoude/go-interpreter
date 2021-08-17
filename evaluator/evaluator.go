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
	case *ast.BlockStatement: // {}
		return evalBlockStatements(node.Statements)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		return &object.ReturnValue{Value: val}
		// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	switch {
	case isTruthy(condition):
		return Eval(ie.Consequence)
	case ie.Alternative != nil:
		return Eval(ie.Alternative)
	default:
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	// “truthy” means: it’s not null and it’s not false
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ: // 左右都是integer数据类型直接进行运算
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ: // 左右都是Boolean数据类型
		return evalBooleanInfixExpression(operator, left, right)
	}
	return NULL
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	switch operator {
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
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

/*
需要考虑, 应该返回 10, if 嵌套返回值的问题
if (10 > 1) {
if (10 > 1) {
return 10;
}
return 1;
}
*/
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s) // 解析最后一条语句才是返回值

		if returnValue, ok := result.(*object.ReturnValue); ok {
			// program 外层遇到 return 返回
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s)
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
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
