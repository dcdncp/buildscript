package std

import (
	"bscript/parser/ast"
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
)

type Function struct {
	*object.Object
	Env      *value.Env
	Self     value.Value
	Params   []string
	Body     ast.Stmt
	Variadic bool
	static   bool
}

func EmptyFunction() value.Value {
	return &Function{object.EmptyObject().(*object.Object), nil, nil, make([]string, 0), nil, false, false}
}

var FunctionType *types.Type

func NewFunction(env *value.Env, s value.Value, params []string, body ast.Stmt, variadic bool, static bool) *Function {
	self := FunctionType.New().(*Function)
	self.Env = env
	self.Self = s
	self.Params = params
	self.Body = body
	self.Variadic = variadic
	return self
}
func NewMethod(env *value.Env, params []string, body ast.Stmt, variadic bool) *Function {
	self := FunctionType.New().(*Function)
	self.Env = env
	self.Self = nil
	self.Params = params
	self.Body = body
	self.Variadic = variadic
	self.static = false
	return self
}
func (o *Function) Call(env *value.Env, args ...value.Value) (value.Value, state.State, bool) {
	e := o.Env.NewChild()
	v, stt := CheckArgsCount(env, len(args), len(o.Params), o.Variadic)
	if stt.IsNotOkay() {
		return v, stt, true
	}
	if o.Variadic {
		last := o.Params[len(o.Params)-1]
		index := 0
		for i, param := range o.Params {
			index = i
			if param == last {
				break
			}
			v := args[i]
			if !e.Const(param, v) {
				v, stt := ThrowException(e, fmt.Sprintf("parameter '%s' has already existed", param))
				return v, stt, true
			}
		}
		// if index != 0 {
		// 	index += 1
		// }
		tail := NewTuple(args[index:])
		if !e.Const(last, tail) {
			v, stt := ThrowException(e, fmt.Sprintf("parameter '%s' has already existed", last))
			return v, stt, true
		}
	} else {
		for i, param := range o.Params {
			v := args[i]
			if !e.Const(param, v) {
				v, stt := ThrowException(e, fmt.Sprintf("parameter '%s' has already existed", param))
				return v, stt, true
			}
		}
	}
	v, stt = e.Global.EvalStmt(e, o.Body)
	if stt == state.Return {
		return v, state.Ok, true
	} else if stt == state.Error {
		return v, stt, true
	} else if stt == state.Break {
		v, stt := ThrowException(e, "'break' statement is out of the loop")
		return v, stt, true
	} else if stt == state.Continue {
		v, stt := ThrowException(e, "'continue' statement is out of the loop")
		return v, stt, true
	}
	return NewNull(), state.Ok, true
}
func (o *Function) IsMethod() bool {
	return !o.static
}
func (o *Function) IsFunction() bool {
	return true
}
func (o *Function) BindFunction(self value.Value) bool {
	o.Self = self
	return true
}
func (o *Function) Copy() value.Value {
	return NewFunction(o.Env, o.Self, o.Params, o.Body, o.Variadic, o.static)
}