package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"slices"
	"strings"
)

type Array struct {
	*object.Object
	Values []value.Value
}

func EmptyArray() value.Value {
	return &Array{object.EmptyObject().(*object.Object), make([]value.Value, 0)}
}

var ArrayType = types.NewAtomType("Array", EmptyArray)

func NewArray(values []value.Value) *Array {
	self := ArrayType.New().(*Array)
	self.Values = values
	return self
}
func (o *Array) Copy() value.Value {
	return NewArray(o.Values)
}

func arrayString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Array)
	strs := make([]string, 0)
	for _, i := range self.Values {
		v, stt := ToString(env, i)
		if stt.IsNotOkay() {
			return v, stt
		}
		strs = append(strs, v.(*String).Value)
	}
	return NewString("[" + strings.Join(strs, ", ") + "]"), state.Ok
}
func arrayIter(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Array)
	return NewStaticExtern(func(env *value.Env, args []value.Value) (value.Value, state.State) {
		yield := args[0]
		for _, i := range self.Values {
			v, stt := Call(env, yield, i)
			if stt == state.Continue {
				continue
			} else if stt == state.Break {
				break
			} else if stt.IsNotOkay() {
				return v, stt
			}
		}
		return nil, state.Ok
	}), state.Ok
}
func arrayIndex(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Array)
	index := args[1].(*Int).Value
	return self.Values[index], state.Ok
}
func arrayPush(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Array)
	self.Values = append(self.Values, args[1:]...)
	return nil, state.Ok
}
func arrayInsert(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Array)
	at := args[1].(*Int).Value
	values := args[2:]
	s := self.Values
	self.Values = slices.Insert(s, int(at), values...)
	return nil, state.Ok
}
