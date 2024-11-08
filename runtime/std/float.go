package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
)

type Float struct {
	*object.Object
	Value float64
}

func EmptyFloat() value.Value {
	return &Float{object.EmptyObject().(*object.Object), 0.0}
}

var FloatType *types.Type

func NewFloat(value float64) *Float {
	self := FloatType.New().(*Float)
	self.Value = value
	return self
}
func (o *Float) Copy() value.Value {
	return NewFloat(o.Value)
}

func floatString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*Float)
	return NewString(fmt.Sprint(self.Value)), state.Ok
}
