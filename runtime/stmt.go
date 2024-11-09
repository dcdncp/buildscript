package runtime

import (
	"bscript/parser/ast"
	"bscript/runtime/state"
	"bscript/runtime/std"
	"bscript/runtime/value"
	"fmt"
)

func IsTrue(v value.Value) bool {
	if v.Type() == std.BoolType {
		return v.(*std.Bool).Value
	} else if v.Type() == std.IntType {
		return v.(*std.Int).Value != 0
	} else if v.Type() == std.FloatType {
		return v.(*std.Float).Value != 0
	} else if v.Type() == std.NullType {
		return false
	}
	return true
}
func EvalProgramStmt(env *value.Env, stmt *ast.ProgramStmt) (value.Value, state.State) {
	for _, s := range stmt.Body {
		v, stt := EvalStmt(env, s)
		if stt.IsNotOkay() {
			return v, stt
		}
	}
	return nil, state.Ok
}
func EvalStmt(env *value.Env, stmt ast.Stmt) (value.Value, state.State) {
	ps := env.Global.Node
	defer env.SetCurrentNode(ps)
	env.SetCurrentNode(stmt)

	switch stmt.Kind() {
	case ast.ProgramStmtKind:
		return EvalProgramStmt(env, stmt.(*ast.ProgramStmt))
	case ast.ConstStmtKind:
		return EvalConstStmt(env, stmt.(*ast.ConstStmt))
	case ast.DefineStmtKind:
		return EvalDefineStmt(env, stmt.(*ast.DefineStmt))
	case ast.BlockStmtKind:
		return EvalBlockStmt(env, stmt.(*ast.BlockStmt))
	case ast.IfStmtKind:
		return EvalIfStmt(env, stmt.(*ast.IfStmt))
	case ast.ForStmtKind:
		return EvalForStmt(env, stmt.(*ast.ForStmt))
	case ast.WhileStmtKind:
		return EvalWhileStmt(env, stmt.(*ast.WhileStmt))
	case ast.ExprStmtKind:
		return EvalExpr(env, stmt.(*ast.ExprStmt).Expr)
	case ast.BreakStmtKind:
		return nil, state.Break
	case ast.ContinueStmtKind:
		return nil, state.Continue
	case ast.ReturnStmtKind:
		s := stmt.(*ast.ReturnStmt)
		if s.Value != nil {
			v, stt := EvalExpr(env, s.Value)
			if stt.IsNotOkay() {
				return v, stt
			}
			return v, state.Return
		}
		return std.NewNull(), state.Return
	case ast.ThrowStmtKind:
		s := stmt.(*ast.ThrowStmt)
		v, stt := EvalExpr(env, s.Value)
		if stt.IsNotOkay() {
			return v, stt
		}
		return v, state.Error
	}
	panic("unreachable")
}
func EvalConstStmt(env *value.Env, stmt *ast.ConstStmt) (value.Value, state.State) {
	if stmt.Id.Kind() != ast.IdentExprKind {
		return std.ThrowException(env, fmt.Sprintf("can not use '%s' as ident", stmt.Id))
	}
	id := stmt.Id.(*ast.IdentExpr).Value
	v, stt := EvalExpr(env, stmt.Value)
	if stt.IsNotOkay() {
		return v, stt
	}
	if !env.Const(id, v) {
		return std.ThrowException(env, fmt.Sprintf("variable '%s' is already exists", id))
	}
	return nil, state.Ok
}
func EvalDefineStmt(env *value.Env, stmt *ast.DefineStmt) (value.Value, state.State) {
	ident := stmt.Ident.Value
	params := make([]string, 0)
	for _, p := range stmt.Params {
		params = append(params, p.Value)
	}
	f := std.NewFunction(env, nil, params, stmt.Body, stmt.Variadic, true)
	if !env.Const(ident, f) {
		return std.ThrowException(env, fmt.Sprintf("'%s' is already exist", ident))
	}
	return nil, state.Ok
}
func EvalBlockStmt(env *value.Env, stmt *ast.BlockStmt) (value.Value, state.State) {
	e := env.NewChild()
	for _, s := range stmt.Body {
		v, stt := EvalStmt(e, s)
		if stt.IsNotOkay() {
			return v, stt
		}
	}
	return nil, state.Ok
}
func EvalIfStmt(env *value.Env, stmt *ast.IfStmt) (value.Value, state.State) {
	condtion, stt := EvalExpr(env, stmt.Condition)
	if stt.IsNotOkay() {
		return condtion, stt
	}
	if IsTrue(condtion) {
		return EvalStmt(env, stmt.Primary)
	} else if stmt.Secondary != nil {
		return EvalStmt(env, stmt.Secondary)
	}
	return nil, state.Ok
}

// func[I](yield func(I)bool)
func EvalForStmt(env *value.Env, stmt *ast.ForStmt) (value.Value, state.State) {
	in, stt := EvalExpr(env, stmt.In)
	if stt.IsNotOkay() {
		return in, stt
	}
	ln := env.Global.Node
	env.SetCurrentNode(stmt.In)
	var iterator value.Value
	if in.Type() == std.FunctionType {
		iterator = in
	} else {
		iterator, stt = std.GetIter(env, in)
		if stt.IsNotOkay() {
			return iterator, stt
		}
	}
	env.SetCurrentNode(ln)
	if stt.IsNotOkay() {
		return iterator, stt
	}
	yield := std.NewStaticExtern(func(e *value.Env, args []value.Value) (value.Value, state.State) {
		env := e.NewChild()
		env.Const(stmt.Ident.Value, args[0])
		return EvalStmt(env, stmt.Body)
	})
	return std.Call(env, iterator, yield)
}
func EvalWhileStmt(env *value.Env, stmt *ast.WhileStmt) (value.Value, state.State) {
	for {
		condition, stt := EvalExpr(env, stmt.Condition)
		if stt.IsNotOkay() {
			return condition, stt
		}
		if !IsTrue(condition) {
			break
		}
		v, stt := EvalStmt(env, stmt.Body)
		if stt == state.Break {
			break
		} else if stt == state.Continue {
			continue
		} else if stt.IsNotOkay() {
			return v, stt
		}
	}
	return nil, state.Ok
}
