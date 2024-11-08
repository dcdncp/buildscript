package std

import (
	"bscript/runtime/object"
	"bscript/runtime/types"
	"bscript/runtime/value"
)

type Module struct {
	*object.Object
	Env *value.Env
}

func EmptyModule() value.Value {
	return &Module{object.EmptyObject().(*object.Object), nil}
}

var ModuleType *types.Type

func NewModule(env *value.Env) *Module {
	self := ModuleType.New().(*Module)
	self.Env = env
	return self
}
func (o *Module) GetField(name string) (value.Value, bool) {
	item, exists := o.Env.GetVariables()[name]
	if !exists {
		return nil, false
	}
	return item.Value, true
}
func (o *Module) SetField(name string, v value.Value) bool {
	item, exists := o.Env.GetVariables()[name]
	if !exists {
		return false
	}
	if item.Const {
		return false
	}
	item.Value = v
	return true
}
func (o *Module) Copy() value.Value {
	return NewModule(o.Env)
}