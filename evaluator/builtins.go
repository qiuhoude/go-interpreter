package evaluator

import "github.com/qiuhoude/go-interpreter/object"

var builtins = map[string]object.Object{
	"len":   makeBuiltin(builtinLen),
	"first": makeBuiltin(builtinFirst),
	"last":  makeBuiltin(builtinLast),
	"rest":  makeBuiltin(builtinRest),
	"push":  makeBuiltin(builtinPush),
}

func makeBuiltin(fn object.BuiltinFunction) *object.Builtin {
	return &object.Builtin{Fn: fn}
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}

	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func arrayOp(op func(*object.Array) object.Object, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `array operate` must be ARRAY, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	return op(arr)
}

func builtinFirst(args ...object.Object) object.Object {
	return arrayOp(func(arr *object.Array) object.Object {
		if len(arr.Elements) == 0 {
			return NULL
		}
		return arr.Elements[0]
	}, args...)
}

func builtinLast(args ...object.Object) object.Object {
	return arrayOp(func(arr *object.Array) object.Object {
		if len(arr.Elements) == 0 {
			return NULL
		}
		return arr.Elements[len(arr.Elements)-1]
	}, args...)
}

func builtinRest(args ...object.Object) object.Object {
	return arrayOp(func(arr *object.Array) object.Object {
		if len(arr.Elements) == 0 {
			return NULL
		}
		length := len(arr.Elements)
		newElements := make([]object.Object, length-1, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}, args...)
}

func builtinPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s",
			args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length+1, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]
	return &object.Array{Elements: newElements}
}
