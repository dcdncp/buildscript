package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
)

type String struct {
	*object.Object
	Value string
}

func EmptyString() value.Value {
	return &String{object.EmptyObject().(*object.Object), ""}
}

var StringType *types.Type

func NewString(value string) *String {
	self := StringType.New().(*String)
	self.Value = value
	return self
}
func (o *String) Copy() value.Value {
	return NewString(o.Value)
}

func stringAdd(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*String)
	other := args[1].(*String)
	return NewString(self.Value + other.Value), state.Ok
}
func stringString(env *value.Env, args []value.Value) (value.Value, state.State) {
	return NewString("\"" + args[0].(*String).Value + "\""), state.Ok
}
