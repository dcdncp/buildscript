package object

import (
	"bscript/runtime/state"
	"bscript/runtime/symbol"
	"bscript/runtime/value"
)

type Object struct {
	parent value.Value
	typ    value.Value
	fields value.VariableMap
	locked bool
}

func EmptyObject() value.Value {
	return &Object{nil, nil, make(value.VariableMap), true}
}

func (o *Object) callMethod(e *value.Env, name string, args ...value.Value) (value.Value, state.State, bool) {
	method, exists := o.GetField(name)
	if !exists {
		return nil, state.Error, false
	}
	v, stt, exists := method.Call(e, args...)
	if !exists {
		return nil, state.Error, false
	}
	return v, stt, true
}

func (o *Object) GetFields() value.VariableMap {
	return o.fields
}
func (o *Object) GetField(name string) (value.Value, bool) {
	item, exists := o.fields[name]
	if !exists {
		if o.parent != nil {
			return o.parent.GetField(name)
		}
		return nil, false
	}
	return item.Value, true
}
func (o *Object) SetField(name string, v value.Value) bool {
	item, exists := o.fields[name]
	if !exists {
		if o.parent != nil {
			if !o.parent.SetField(name, v) {
				if o.locked {
					return false
				}
				o.fields[name] = &value.Variable{Value: v, Const: false}
			}
			return true
		}
		if !o.locked {
			o.fields[name] = &value.Variable{Value: v, Const: false}
			return true
		}
		return false
	}
	if item.Const {
		return false
	}
	item.Value = v
	return true
}
func (o *Object) ConstField(name string, v value.Value) bool {
	_, exists := o.fields[name]
	if exists {
		return false
	}
	o.fields[name] = &value.Variable{Value: v, Const: true}
	return true
}
func (o *Object) GetParent() value.Value {
	return o.parent
}
func (o *Object) SetParent(value value.Value) {
	o.parent = value
}
func (o *Object) GetKey(e *value.Env, key value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.IndexGet, key)
}
func (o *Object) SetKey(e *value.Env, key value.Value, value value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.IndexSet, key, value)
}
func (o *Object) Add(e *value.Env, other value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Add, other)
}
func (o *Object) Sub(e *value.Env, other value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Sub, other)
}
func (o *Object) Mul(e *value.Env, other value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Mul, other)
}
func (o *Object) Div(e *value.Env, other value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Div, other)
}
func (o *Object) Call(e *value.Env, args ...value.Value) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Call, args...)
}
func (o *Object) String(e *value.Env) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.String)
}
func (o *Object) GetIter(e *value.Env) (value.Value, state.State, bool) {
	return o.callMethod(e, symbol.Iter)
}
func (o *Object) IsFunction() bool {
	return false
}
func (o *Object) IsMethod() bool {
	return false
}
func (o *Object) BindFunction(self value.Value) bool {
	return false
}
func (o *Object) Copy() value.Value {
	no := EmptyObject().(*Object)
	no.typ = o.typ
	if o.parent != nil {
		no.parent = o.parent.Copy()
	}
	no.fields = o.fields
	no.locked = o.locked
	return no
}
func (o *Object) Type() value.Value {
	// if o.typ == nil {
	// 	return TypeType
	// }
	return o.typ
}
func (o *Object) Lock() {
	o.locked = true
}
func (o *Object) Unlock() {
	o.locked = false
}
func (o *Object) SetType(typ value.Value) {
	o.typ = typ
}
