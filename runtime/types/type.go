package types

import (
	"bscript/runtime/value"
	"bscript/runtime/object"
)

type Type struct {
	*object.Object
	Default func() value.Value
	Name    string
	Proto   value.Value
}

func EmptyType() value.Value {
	return &Type{object.EmptyObject().(*object.Object), nil, "", nil}
}

var TypeType = NewRawType("Type", EmptyType)
var ProtoType = NewRawType("Proto", object.EmptyObject)

func NewRawType(name string, defaultF func() value.Value) *Type {
	return &Type{object.EmptyObject().(*object.Object), defaultF, name, object.EmptyObject().(*object.Object)}
}
func NewType(name string, parent *Type) *Type {
	self := TypeType.New().(*Type)
	if parent != nil {
		self.SetParent(parent)
	}
	self.Name = name
	self.Default = object.EmptyObject
	proto := ProtoType.New()
	self.Proto = proto
	return self
}
func NewAtomType(name string, defaultF func() value.Value) *Type {
	self := TypeType.New().(*Type)
	self.Name = name
	self.Default = defaultF
	proto := ProtoType.New()
	self.Proto = proto
	return self
}
func (o *Type) Type() value.Value {
	return TypeType
}
func (o *Type) Copy() value.Value {
	return o
}
func (t *Type) init(self value.Value) value.Value {
	self.Unlock()
	self.SetType(t)
	parent := t.GetParent()
	if parent != nil {
		ptyp, ok := parent.(*Type)
		if ok && ptyp != nil {
			parentObj := ptyp.Default()
			ptyp.init(parentObj)
			self.SetParent(parentObj)
		}
	}
	if t.Proto != nil {
		for name, item := range t.Proto.GetFields() {
			field := item.Value.Copy()
			if field.IsMethod() {
				field.BindFunction(self)
				self.ConstField(name, field)
			} else {
				self.SetField(name, field)
			}
		}
	}
	self.Lock()
	return self
}
func (t *Type) New() value.Value {
	self := t.Default()
	t.init(self)
	return self
}
