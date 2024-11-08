package std

import (
	"bscript/errors"

	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"

	"fmt"
)

type Exception struct {
	*object.Object
	Errors    []errors.Error
	Callstack []string
}

func EmptyException() value.Value {
	return &Exception{object.EmptyObject().(*object.Object), make([]errors.Error, 0), make([]string, 0)}
}

var ExceptionType *types.Type

func NewException(errors []errors.Error, callstack []string) *Exception {
	self := ExceptionType.New().(*Exception)
	self.Errors = errors
	self.Callstack = callstack
	return self
}
func ThrowException(env *value.Env, msg string) (value.Value, state.State) {
	n := env.GetCurrentNode()
	err := errors.NewRuntimeError(
		env.Global.SourceFile, env.Global.Source, msg, n.Start(), n.End())
	return NewException([]errors.Error{err}, env.Global.Callstack), state.Error
}
func (o *Exception) Copy() value.Value {
	return NewException(o.Errors, o.Callstack)
}

func exceptionInit(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Exception)
	msg := args[1].(*String).Value
	n := env.GetCurrentNode()
	fmt.Println(n, n.Start().Col, n.End().Col)
	err := errors.NewRuntimeError(env.Global.SourceFile, env.Global.Source, msg, n.Start(), n.End())
	self.Errors = []errors.Error{err}
	self.Callstack = env.Global.Callstack
	return nil, state.Ok
}
func exceptionString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Exception)
	return NewString(fmt.Sprintf("<exception %p>", self)), state.Ok
}
