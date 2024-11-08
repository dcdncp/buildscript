package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
)

type Int struct {
	*object.Object
	Value int64
}

func EmptyInt() value.Value {
	return &Int{object.EmptyObject().(*object.Object), 0}
}

var IntType *types.Type

func NewInt(value int64) *Int {
	self := IntType.New().(*Int)
	self.Value = value
	return self
}
func (o *Int) Copy() value.Value {
	return NewInt(o.Value)
}

func intAdd(env *value.Env, args []value.Value) (value.Value, state.State) {
	a := args[0].(*Int)
	b := args[1].(*Int)
	return NewInt(a.Value + b.Value), state.Ok
}
func intSub(env *value.Env, args []value.Value) (value.Value, state.State) {
	a := args[0].(*Int)
	b := args[1].(*Int)
	return NewInt(a.Value - b.Value), state.Ok
}
func intMul(env *value.Env, args []value.Value) (value.Value, state.State) {
	a := args[0].(*Int)
	b := args[1].(*Int)
	return NewInt(a.Value * b.Value), state.Ok
}
func intDiv(env *value.Env, args []value.Value) (value.Value, state.State) {
	a := args[0].(*Int)
	b := args[1].(*Int)
	return NewFloat(float64(a.Value) / float64(b.Value)), state.Ok
}
func intString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Int)
	return NewString(fmt.Sprint(self.Value)), state.Ok
}
