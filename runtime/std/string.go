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

func stringInit(env *value.Env, args []value.Value) (value.Value, state.State) {
	v, stt := CheckArgsCount(env, len(args), 2, false)
	if stt.IsNotOkay() {
		return v, stt
	}

	v, stt = CheckType(env, args[0], StringType)
	if stt.IsNotOkay() {
		return v, stt
	}
	self := v.(*String)
	
	value := args[1]
	v, stt = ToString(env, value)
	if stt.IsNotOkay() {
		return v, stt
	}

	self.Value = v.(*String).Value

	return nil, state.Ok
}
func stringAdd(env *value.Env, args []value.Value) (value.Value, state.State) {
	v, stt := CheckArgsCount(env, len(args), 2, false)
	if stt.IsNotOkay() {
		return v, stt
	}

	v, stt = CheckType(env, args[0], StringType)
	if stt.IsNotOkay() {
		return v, stt
	}
	self := v.(*String)

	v, stt = CheckType(env, args[1], StringType)
	if stt.IsNotOkay() {
		return v, stt
	}
	other := v.(*String)

	return NewString(self.Value + other.Value), state.Ok
}
func stringString(env *value.Env, args []value.Value) (value.Value, state.State) {
	return NewString("\"" + args[0].(*String).Value + "\""), state.Ok
}
