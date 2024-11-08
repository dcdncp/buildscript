package fs

import (
	"bscript/runtime/object"
	"bscript/runtime/state"
	"bscript/runtime/std"
	"bscript/runtime/types"
	"bscript/runtime/value"
	"fmt"
	"os"
)

type FileS struct {
	*object.Object
	file os.File
}

var File *types.Type

func EmptyFile() value.Value {
	return &FileS{object.EmptyObject().(*object.Object), os.File{}}
}

func NewFile(file os.File) *FileS {
	self := File.New().(*FileS)
	self.file = file
	return self
}
func (f *FileS) Copy() value.Value {
	return NewFile(f.file)
}

func fileString(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*FileS)
	return std.NewString(fmt.Sprintf("<file %p>", self)), state.Ok
}
func fileClose(env *value.Env, args []value.Value) (value.Value, state.State) {
	self := args[0].(*FileS)
	err := self.file.Close()
	if err != nil {
		return std.ThrowException(env, err.Error())
	}
	return nil, state.Ok
}
