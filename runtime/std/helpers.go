package std

import (
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
)

func CheckArgsCount(env *value.Env, count, expected int, variadic bool) (value.Value, state.State) {
	if variadic {
		if count < expected-1 {
			return ThrowException(env, "too few arguments")
		}
	} else {
		if count < expected {
			return ThrowException(env, "too few arguments")
		} else if count > expected {
			return ThrowException(env, "too many arguments")
		}
	}
	return nil, state.Ok
}
func CheckType(env *value.Env, v value.Value, typ value.Value) (value.Value, state.State) {
	if v.Type().(*types.Type) != typ.(*types.Type) {
		return ThrowException(env, fmt.Sprintf("value has incorrect type, '%s' expected", typ.(*types.Type).Name))
	}
	return v, state.Ok
}
func CheckTypes(env *value.Env, values []value.Value, typ value.Value) (value.Value, state.State) {
	for _, v := range values {
		if v.Type().(*types.Type) != typ.(*types.Type) {
			return ThrowException(env, fmt.Sprintf("value has incorrect type, '%s' expected", typ.(*types.Type).Name))
		}
	}
	return NewTuple(values), state.Ok
}
