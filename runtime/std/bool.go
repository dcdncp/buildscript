package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
)

type Bool struct {
	*object.Object
	Value bool
}

func EmptyBool() value.Value {
	return &Bool{object.EmptyObject().(*object.Object), false}
}

var BoolType *types.Type

func NewBool(value bool) *Bool {
	self := BoolType.New().(*Bool)
	self.Value = value
	return self
}
func (o *Bool) Copy() value.Value {
	return NewBool(o.Value)
}

func boolString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Bool)
	return NewString(fmt.Sprint(self.Value)), state.Ok
}
