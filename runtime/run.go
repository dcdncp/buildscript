package runtime

import (
	"bscript/errors"
	"bscript/parser"
	"bscript/parser/ast"
	"bscript/runtime/state"
	"bscript/runtime/std"
	"path"
)

func RunSource(source string) []errors.Error {
	prog, errs := parser.ParseSource(source)
	if len(errs) > 0 {
		return errs
	}
	v, stt := EvalProgram(prog, "<virtual>", source)
	if stt.IsNotOkay() {
		return v.(*std.Exception).Errors
	}
	return make([]errors.Error, 0)
}

func RunProgram(prog *ast.ProgramStmt, file, source string) []errors.Error {
	e := std.Env.NewChild()
	e.Global.Source = source
	e.Global.SourceFile = file
	e.Global.WorkingDir = path.Dir(file)
	v, stt := EvalProgramStmt(e, prog)
	if stt == state.Error {
		return v.(*std.Exception).Errors
	} else if stt == state.Return {
		e, _ := std.ThrowException(e, "'return' statement is out of the function")
		return e.(*std.Exception).Errors
	} else if stt == state.Continue {
		e, _ := std.ThrowException(e, "'continue' statement is out of the loop")
		return e.(*std.Exception).Errors
	} else if stt == state.Break {
		e, _ := std.ThrowException(e, "'break' statement is out of the loop")
		return e.(*std.Exception).Errors
	}
	return make([]errors.Error, 0)
}

func RunFile(file string) []errors.Error {
	prog, errs, source := parser.ParseFile(file)
	if len(errs) > 0 {
		return errs
	}
	return  RunProgram(prog, file, source)
}
