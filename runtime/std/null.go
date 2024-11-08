package std

import (
	"bscript/runtime/types"
	"bscript/runtime/value"
)

var NullType *types.Type

func NewNull() value.Value {
	return NullType.New()
}
