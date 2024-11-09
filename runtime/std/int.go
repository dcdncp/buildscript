package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
	"strconv"
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

func intInit(env *value.Env, args []value.Value) (value.Value, state.State) {
	v, stt := CheckArgsCount(env, len(args), 2, false)
	if stt.IsNotOkay() {
		return v, stt
	}

	v, stt = CheckType(env, args[0], IntType)
	if stt.IsNotOkay() {
		return v, stt
	}
	self := v.(*Int)

	v = args[1]
	if v.Type() == IntType {
		self.Value = v.(*Int).Value
	} else if v.Type() == StringType {
		value, err := strconv.ParseInt(v.(*String).Value, 10, 64)
		if err != nil {
			return ThrowException(env, err.Error())
		}
		self.Value = value
	} else {
		return ThrowException(env, "value has incorrect type, expected Int or String")
	}

	return nil, state.Ok
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
