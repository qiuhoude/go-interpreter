package evaluator

import (
	"fmt"
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
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

		// expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}
	return nil
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
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
	case left.Type() == object.INTEGER_OBJ && right.Type() == left.Type(): // 左右都是integer数据类型直接进行运算
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == left.Type(): // 左右都是Boolean数据类型
		return evalBooleanInfixExpression(operator, left, right)
	case right.Type() != left.Type(): // 左右两边类型不相等
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())

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
		return newError("unknown operator: %s%s", operator, right.Type())

	}
}

func evalMinusOrPlusOperatorExpression(op token.TokenType, right object.Object) object.Object {
	r, ok := right.(*object.Integer)
	if !ok {
		return newError("unknown operator: -%s", right.Type())
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
		switch result := result.(type) {
		case *object.ReturnValue:
			// 此处运用于只有一层return语句时有效,套会导致只有最外层的return语句有效
			return result.Value
		case *object.Error: // 有错误提前返回
			return result
		}
	}
	return result
}

func evalBlockStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, s := range stmts {
		result = Eval(s)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				// 返回return本身, 表示外层也是获得statement的object也是return,不往下继续进行解析到此结束
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
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
