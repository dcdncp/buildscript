package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"strings"
)

type Tuple struct {
	*object.Object
	Values []value.Value
}

func EmptyTuple() value.Value {
	return &Tuple{object.EmptyObject().(*object.Object), make([]value.Value, 0)}
}

var TupleType *types.Type

func NewTuple(values []value.Value) *Tuple {
	self := TupleType.New().(*Tuple)
	self.Values = values
	return self
}
func (o *Tuple) Copy() value.Value {
	return NewTuple(o.Values)
}
func tupleInit(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Tuple)
	self.Values = args[1:]
	return nil, state.Ok
}
func tupleString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Tuple)
	strs := make([]string, 0)
	for _, i := range self.Values {
		v, stt := ToString(env, i)
		if stt.IsNotOkay() {
			return v, stt
		}
		strs = append(strs, v.(*String).Value)
	}
	return NewString("(" + strings.Join(strs, ", ") + ")"), state.Ok
}
func tupleIndexSet(env *value.Env, args []value.Value) (value.Value, state.State) {
	v, stt := CheckArgsCount(env, len(args), 3, false)
	if stt.IsNotOkay() {
		return v, stt
	}

	v, stt = CheckType(env, args[0], TupleType)
	if stt.IsNotOkay() {
		return v, stt
	}
	self := v.(*Tuple)

	v, stt = CheckType(env, args[1], IntType)
	if stt.IsNotOkay() {
		return v, stt
	}
	index := v.(*Int).Value

	value := args[2]

	self.Values[index] = value
	return nil, state.Ok
}
func tupleIndexGet(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Tuple)
	index := args[1].(*Int).Value
	return self.Values[index], state.Ok
}
func tupleIter(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Tuple)
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
