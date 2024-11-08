package fs

import (
	"bscript/runtime/std"
	"bscript/runtime/symbol"
	"bscript/runtime/types"
	"bscript/runtime/value"
)

var Env *value.Env

func AddValue(name string, v value.Value) {
	Env.Const(name, v)
}
func AddFunction(name string, f std.ExternFunc) {
	Env.Const(name, std.NewStaticExtern(f))
}
func AddTypeMethod(typ *types.Type, name string, f std.ExternFunc) {
	typ.Proto.ConstField(name, std.NewExternMethod(f))
}
func AddMethod(v value.Value, name string, f std.ExternFunc) {
	v.ConstField(name, std.NewExtern(v, f, false))
}

func init() {
	Env = value.NewEnv()

	File = types.NewAtomType("File", EmptyFile)

	AddTypeMethod(File, symbol.String, fileString)
	AddTypeMethod(File, "close", fileClose)

	AddFunction("read_file", readFile)
	AddFunction("write_file", writeFile)
	
	std.Modules["fs"] = Env
}
