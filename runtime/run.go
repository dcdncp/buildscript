package runtime

import (
	"bscript/errors"
	"bscript/parser/ast"
	"bscript/runtime/std"
)

func RunProgram(prog *ast.ProgramStmt, file, source string) []errors.Error  {
	v, stt := EvalProgram(prog, file, source)
	if stt.IsNotOkay() {
		return v.(*std.Exception).Errors
	}
	return make([]errors.Error, 0)
}
func RunSource(source string) []errors.Error {
	v, stt := EvalSource(source)
	if stt.IsNotOkay() {
		return v.(*std.Exception).Errors
	}
	return make([]errors.Error, 0)
}
func RunFile(file string) []errors.Error {
	v, stt := EvalFile(file)
	if stt.IsNotOkay() {
		return v.(*std.Exception).Errors
	}
	return make([]errors.Error, 0)
}
