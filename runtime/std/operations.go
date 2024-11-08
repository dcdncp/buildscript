package std

import (
	"bscript/runtime/state"
	"bscript/runtime/value"
	"fmt"
)

func Add(env *value.Env, self value.Value, other value.Value) (value.Value, state.State) {
	v, stt, exists := self.Add(env, other)
	if !exists {
		return ThrowException(env, "object has not '+' operator")
	}
	return v, stt
}
func Sub(env *value.Env, self value.Value, other value.Value) (value.Value, state.State) {
	v, stt, exists := self.Sub(env, other)
	if !exists {
		return ThrowException(env, "object has not '-' operator")
	}
	return v, stt
}
func Mul(env *value.Env, self value.Value, other value.Value) (value.Value, state.State) {
	v, stt, exists := self.Mul(env, other)
	if !exists {
		return ThrowException(env, "object has not '*' operator")
	}
	return v, stt
}
func Div(env *value.Env, self value.Value, other value.Value) (value.Value, state.State) {
	v, stt, exists := self.Div(env, other)
	if !exists {
		return ThrowException(env, "object has not '/' operator")
	}
	return v, stt
}
func GetKey(env *value.Env, self value.Value, key value.Value) (value.Value, state.State) {
	v, stt, exists := self.GetKey(env, key)
	if !exists {
		return ThrowException(env, "object can not be indexed")
	}
	return v, stt
}
func SetKey(env *value.Env, self value.Value, key value.Value, v value.Value) (value.Value, state.State) {
	v, stt, exists := self.SetKey(env, key, v)
	if !exists {
		return ThrowException(env, "object can not be indexed")
	}
	return v, stt
}
func ToString(env *value.Env, self value.Value) (value.Value, state.State) {
	v, stt, exists := self.String(env)
	if !exists {
		return NewString(fmt.Sprintf("<object %p>", self)), state.Ok
	}
	return v, stt
}
func Call(env *value.Env, self value.Value, args ...value.Value) (value.Value, state.State) {
	v, stt, exists := self.Call(env, args...)
	if !exists {
		return ThrowException(env, "object is not callable")
	}
	return v, stt
}
func GetIter(env *value.Env, self value.Value) (value.Value, state.State) {
	v, stt, exists := self.GetIter(env)
	if !exists {
		return ThrowException(env, "object is not iterable")
	}
	return v, stt
}
func ForEach(env *value.Env, self value.Value, f func(el value.Value)(value.Value, state.State)) (value.Value, state.State) {
	iterator, stt := GetIter(env, self)
	if stt.IsNotOkay() {
		return iterator, stt
	}
	yield := NewStaticExtern(func(e *value.Env, args []value.Value) (value.Value, state.State) {
		return f(args[0])
	})
	return Call(env, iterator, yield)
}