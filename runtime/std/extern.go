package std

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"bscript/tools"
)

type ExternFunc func(env *value.Env, args []value.Value) (value.Value, state.State)
type Extern struct {
	*object.Object
	Self   value.Value
	Func   ExternFunc
	static bool
}

func EmptyExtern() value.Value {
	return &Extern{object.EmptyObject().(*object.Object), nil, nil, true}
}

var ExternType *types.Type

func NewExtern(s value.Value, f ExternFunc, static bool) *Extern {
	self := ExternType.New().(*Extern)
	self.Self = s
	self.Func = f
	self.static = static
	return self
}
func NewExternMethod(f ExternFunc) *Extern {
	self := ExternType.New().(*Extern)
	self.Func = f
	self.static = false
	return self
}
func NewStaticExtern(f ExternFunc) *Extern {
	self := ExternType.New().(*Extern)
	self.Func = f
	self.static = true
	return self
}
func (o *Extern) IsMethod() bool {
	return !o.static
}
func (o *Extern) IsFunction() bool {
	return true
}
func (o *Extern) BindFunction(self value.Value) bool {
	o.Self = self
	return true
}
func (o *Extern) Call(e *value.Env, args ...value.Value) (value.Value, state.State, bool) {
	if !o.static {
		v, stt := o.Func(e, tools.AppendFront(args, o.Self))
		return v, stt, true
	}
	v, stt := o.Func(e, args)
	return v, stt, true
}
func (o *Extern) Copy() value.Value {
	return NewExtern(o.Self, o.Func, o.static)
}
