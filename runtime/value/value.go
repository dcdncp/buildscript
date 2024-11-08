package value

import (
	"bscript/runtime/state"
)

type Value interface {
	SetParent(parent Value)
	GetParent() Value
	GetFields() VariableMap
	GetField(name string) (Value, bool)
	SetField(name string, value Value) bool
	ConstField(name string, value Value) bool
	GetIter(e *Env) (Value, state.State, bool)
	GetKey(e *Env, key Value) (Value, state.State, bool)
	SetKey(e *Env, key Value, value Value) (Value, state.State, bool)
	Add(e *Env, other Value) (Value, state.State, bool)
	Sub(e *Env, other Value) (Value, state.State, bool)
	Mul(e *Env, other Value) (Value, state.State, bool)
	Div(e *Env, other Value) (Value, state.State, bool)
	Call(e *Env, args ...Value) (Value, state.State, bool)
	String(e *Env) (Value, state.State, bool)
	IsFunction() bool
	IsMethod() bool
	BindFunction(self Value) bool
	Copy() Value
	Type() Value
	Lock()
	Unlock()
	SetType(typ Value)
}
