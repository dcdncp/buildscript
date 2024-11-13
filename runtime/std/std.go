package std

import (
	"bscript/runtime/symbol"
	"bscript/runtime/types"
	"bscript/runtime/value"
)

var Env *value.Env
var Modules map[string]*value.Env = make(map[string]*value.Env)
var LoadedModules map[string]value.Value = make(map[string]value.Value)

func AddValue(name string, v value.Value) {
	Env.Const(name, v)
}
func AddFunction(name string, f ExternFunc) {
	Env.Const(name, NewStaticExtern(f))
}
func AddTypeMethod(typ *types.Type, name string, f ExternFunc) {
	typ.Proto.ConstField(name, NewExternMethod(f))
}
func AddMethod(v value.Value, name string, f ExternFunc) {
	v.ConstField(name, NewExtern(v, f, false))
}

func init() {
	Env = value.NewEnv()

	ExternType = types.NewRawType("Extern", EmptyExtern)
	AddMethod(ExternType, symbol.String, typeString)

	AddMethod(types.TypeType, symbol.String, typeString)
	AddMethod(types.ProtoType, symbol.String, typeString)

	AddTypeMethod(types.TypeType, symbol.String, typeString)
	AddTypeMethod(types.TypeType, symbol.Call, typeCall)

	IntType = types.NewAtomType("Int", EmptyInt)
	FloatType = types.NewAtomType("Float", EmptyFloat)
	BoolType = types.NewAtomType("Bool", EmptyBool)
	StringType = types.NewAtomType("String", EmptyString)
	ArrayType = types.NewAtomType("Array", EmptyArray)
	TupleType = types.NewAtomType("Tuple", EmptyTuple)
	FunctionType = types.NewAtomType("Function", EmptyFunction)
	ExceptionType = types.NewAtomType("Exception", EmptyException)
	ModuleType = types.NewAtomType("Module", EmptyModule)
	NullType = types.NewType("Null", nil)

	IntType.ConstField(symbol.Init, NewStaticExtern(intInit))
	AddTypeMethod(IntType, symbol.Add, intAdd)
	AddTypeMethod(IntType, symbol.Sub, intSub)
	AddTypeMethod(IntType, symbol.Mul, intMul)
	AddTypeMethod(IntType, symbol.Div, intDiv)
	AddTypeMethod(IntType, symbol.String, intString)

	AddTypeMethod(FloatType, symbol.String, floatString)

	AddTypeMethod(BoolType, symbol.String, boolString)

	StringType.ConstField(symbol.Init, NewStaticExtern(stringInit))
	AddTypeMethod(StringType, symbol.Add, stringAdd)
	AddTypeMethod(StringType, symbol.String, stringString)

	TupleType.ConstField(symbol.Init, NewStaticExtern(tupleInit))
	AddTypeMethod(TupleType, symbol.String, tupleString)
	AddTypeMethod(TupleType, symbol.Iter, tupleIter)
	AddTypeMethod(TupleType, symbol.IndexGet, tupleIndexGet)
	AddTypeMethod(TupleType, symbol.IndexSet, tupleIndexSet)

	AddTypeMethod(ArrayType, symbol.String, arrayString)
	AddTypeMethod(ArrayType, symbol.Iter, arrayIter)
	AddTypeMethod(ArrayType, symbol.IndexGet, arrayIndexGet)
	AddTypeMethod(ArrayType, symbol.IndexSet, arrayIndexSet)
	AddTypeMethod(ArrayType, "push", arrayPush)
	AddTypeMethod(ArrayType, "insert", arrayInsert)

	ExceptionType.ConstField(symbol.Init, NewStaticExtern(exceptionInit))
	AddTypeMethod(ExceptionType, symbol.String, exceptionString)

	AddValue("int", IntType)
	AddValue("float", FloatType)
	AddValue("bool", BoolType)
	AddValue("string", StringType)
	AddValue("Tuple", TupleType)
	AddValue("Array", ArrayType)
	AddValue("Exception", ExceptionType)

	AddFunction("import", importModule)
	AddFunction("message", message)
	AddFunction("type_of", typeOf)
}
