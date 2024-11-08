package runtime

import (
	"bscript/parser"
	"bscript/parser/ast"

	"bscript/runtime/state"
	"bscript/runtime/std"
	_ "bscript/runtime/std/fs"
	"bscript/runtime/value"
)

func EvalSource(source string) (value.Value, state.State) {
	prog, errors := parser.ParseSource(source)
	if len(errors) > 0 {
		return std.NewException(errors, make([]string, 0)), state.Error
	}
	return EvalProgram(prog, "<virtual>", source)
}

func EvalProgram(prog *ast.ProgramStmt, file, source string) (value.Value, state.State) {
	e := std.Env.NewChild()
	e.Global.Source = source
	e.Global.SourceFile = file
	v, stt := EvalProgramStmt(e, prog)
	if stt == state.Error {
		return v, stt
	} else if stt == state.Return {
		return std.ThrowException(e, "'return' statement is out of the function")
	} else if stt == state.Continue {
		return std.ThrowException(e, "'continue' statement is out of the loop")
	} else if stt == state.Break {
		return std.ThrowException(e, "'break' statement is out of the loop")
	}
	return std.NewModule(e), state.Ok
}

func EvalFile(file string) (value.Value, state.State) {
	prog, errors, source := parser.ParseFile(file)
	if len(errors) > 0 {
		return std.NewException(errors, make([]string, 0)), state.Error
	}
	return EvalProgram(prog, file, source)
}

func init() {
	std.Env.Global.EvalStmt = EvalStmt
	std.Env.Global.EvalFile = EvalFile
	std.Env.Global.EvalSource = EvalSource
}
