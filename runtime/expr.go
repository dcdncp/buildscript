package runtime

import (
	"bscript/parser/ast"
	"bscript/runtime/state"
	"bscript/runtime/std"
	"bscript/runtime/value"
	"fmt"
)

func EvalExpr(env *value.Env, expr ast.Expr) (value.Value, state.State) {
	pe := env.Global.Node
	defer env.SetCurrentNode(pe)
	env.SetCurrentNode(expr)

	switch expr.Kind() {
	case ast.BinaryExprKind:
		return EvalBinaryExpr(env, expr.(*ast.BinaryExpr))
	case ast.UnaryExprKind:
		return EvalUnaryExpr(env, expr.(*ast.UnaryExpr))
	case ast.IndexExprKind:
		return EvalIndexExpr(env, expr.(*ast.IndexExpr))
	case ast.MemberExprKind:
		return EvalMemberExpr(env, expr.(*ast.MemberExpr))
	case ast.CallExprKind:
		return EvalCallExpr(env, expr.(*ast.CallExpr))
	case ast.AssignExprKind:
		return EvalAssignExpr(env, expr.(*ast.AssignExpr))
	default:
		return EvalPrimaryExpr(env, expr)
	}
}
func EvalAssignExpr(env *value.Env, expr *ast.AssignExpr) (value.Value, state.State) {
	if expr.LHS.Kind() == ast.IdentExprKind {
		ident := expr.LHS.(*ast.IdentExpr)
		value, stt := EvalExpr(env, expr.RHS)
		if stt.IsNotOkay() {
			return value, stt
		}
		if !env.Set(ident.Value, value) {
			return std.ThrowException(env, fmt.Sprintf("'%s' is not assignable", ident.Value))
		}
		return value, state.Ok
	} else if expr.LHS.Kind() == ast.MemberExprKind {
		mexpr := expr.LHS.(*ast.MemberExpr)
		ident := mexpr.Field.Value
		avalue, stt := EvalExpr(env, mexpr.Target)
		if stt.IsNotOkay() {
			return avalue, stt
		}
		value, stt := EvalExpr(env, expr.RHS)
		if stt.IsNotOkay() {
			return value, stt
		}
		if !avalue.SetField(ident, value) {
			return std.ThrowException(env, fmt.Sprintf("'%s' is not assignable", ident))
		}
		return value, state.Ok
	} else if expr.LHS.Kind() == ast.IndexExprKind {
		mexpr := expr.LHS.(*ast.IndexExpr)
		index, stt := EvalExpr(env, mexpr.Index)
		if stt.IsNotOkay() {
			return index, stt
		}
		avalue, stt := EvalExpr(env, mexpr.Target)
		if stt.IsNotOkay() {
			return avalue, stt
		}
		value, stt := EvalExpr(env, expr.RHS)
		if stt.IsNotOkay() {
			return value, stt
		}
		v, stt := std.SetKey(env, avalue, index, value)
		if stt.IsNotOkay() {
			return v, stt
		}
		return value, state.Ok
	} else {
		return std.ThrowException(env, fmt.Sprintf("'%s' is not assignable", expr.LHS))
	}

}
func EvalCallExpr(env *value.Env, expr *ast.CallExpr) (value.Value, state.State) {
	target, stt := EvalExpr(env, expr.Target)
	if stt.IsNotOkay() {
		return target, stt
	}
	args := make([]value.Value, 0)
	for _, e := range expr.Args {
		if IsSpreaded(e) {
			v, stt := EvalExpr(env, e.(*ast.UnaryExpr).Operand)
			if stt.IsNotOkay() {
				return v, stt
			}
			std.ForEach(env, v, func(el value.Value) (value.Value, state.State) {
				args = append(args, el)
				return nil, state.Ok
			})
		} else {
			v, stt := EvalExpr(env, e)
			if stt.IsNotOkay() {
				return v, stt
			}
			args = append(args, v)
		}
	}
	return std.Call(env, target, args...)
}
func EvalIndexExpr(env *value.Env, expr *ast.IndexExpr) (value.Value, state.State) {
	self, stt := EvalExpr(env, expr.Target)
	if stt.IsNotOkay() {
		return self, stt
	}
	index, stt := EvalExpr(env, expr.Index)
	if stt.IsNotOkay() {
		return index, stt
	}
	return std.GetKey(env, self, index)
}
func EvalMemberExpr(env *value.Env, expr *ast.MemberExpr) (value.Value, state.State) {
	v, stt := EvalExpr(env, expr.Target)
	if stt.IsNotOkay() {
		return v, stt
	}
	item, exists := v.GetField(expr.Field.Value)
	if !exists {
		return std.ThrowException(
			env,
			fmt.Sprintf("'%s' does not have field '%s'", expr.Target, expr.Field.Value))
	}
	return item, state.Ok
}
func EvalUnaryExpr(env *value.Env, expr *ast.UnaryExpr) (value.Value, state.State) {
	operator := expr.Operator.Value
	self, state := EvalExpr(env, expr.Operand)
	if state.IsNotOkay() {
		return self, state
	}
	if operator == "..." {
		return std.ThrowException(env, "operator '...' is not available here")
	}
	return std.Add(env, self, nil)
}
func EvalBinaryExpr(env *value.Env, expr *ast.BinaryExpr) (value.Value, state.State) {
	operator := expr.Operator.Value
	self, stt := EvalExpr(env, expr.LHS)
	if stt.IsNotOkay() {
		return self, stt
	}
	o, stt := EvalExpr(env, expr.RHS)
	if stt.IsNotOkay() {
		return o, stt
	}
	var v value.Value
	if operator == "+" {
		v, stt = std.Add(env, self, o)
	} else {
		return std.ThrowException(
			env,
			fmt.Sprintf("'%s' has not '%s' binary operation", expr.LHS, operator))
	}
	return v, stt
}
func IsSpreaded(e ast.Expr) bool {
	return e.Kind() == ast.UnaryExprKind && e.(*ast.UnaryExpr).Operator.Value == "..."
}
func EvalPrimaryExpr(env *value.Env, expr ast.Expr) (value.Value, state.State) {
	switch expr.Kind() {
	case ast.StringExprKind:
		return std.NewString(expr.(*ast.StringExpr).Value), state.Ok
	case ast.IntExprKind:
		return std.NewInt(expr.(*ast.IntExpr).Value), state.Ok
	case ast.FloatExprKind:
		return std.NewFloat(expr.(*ast.FloatExpr).Value), state.Ok
	case ast.BoolExprKind:
		return std.NewBool(expr.(*ast.BoolExpr).Value), state.Ok
	case ast.TupleExprKind:
		exprs := expr.(*ast.TupleExpr).Values
		values := make([]value.Value, 0)
		for _, e := range exprs {
			if IsSpreaded(e) {
				v, stt := EvalExpr(env, e.(*ast.UnaryExpr).Operand)
				if stt.IsNotOkay() {
					return v, stt
				}
				std.ForEach(env, v, func(el value.Value) (value.Value, state.State) {
					values = append(values, el)
					return nil, state.Ok
				})
			} else {
				v, stt := EvalExpr(env, e)
				if stt.IsNotOkay() {
					return v, stt
				}
				values = append(values, v)
			}
		}
		return std.NewTuple(values), state.Ok
	case ast.ArrayExprKind:
		exprs := expr.(*ast.ArrayExpr).Values
		values := make([]value.Value, 0)
		for _, e := range exprs {
			if IsSpreaded(e) {
				v, stt := EvalExpr(env, e.(*ast.UnaryExpr).Operand)
				if stt.IsNotOkay() {
					return v, stt
				}
				std.ForEach(env, v, func(el value.Value) (value.Value, state.State) {
					values = append(values, el)
					return nil, state.Ok
				})
			} else {
				v, stt := EvalExpr(env, e)
				if stt.IsNotOkay() {
					return v, stt
				}
				values = append(values, v)
			}
		}
		return std.NewArray(values), state.Ok
	case ast.DefineExprKind:
		dexpr := expr.(*ast.DefineExpr)
		params := make([]string, 0)
		for _, p := range dexpr.Params {
			params = append(params, p.Value)
		}
		return std.NewFunction(env, nil, params, dexpr.Body, dexpr.Variadic, true), state.Ok
	case ast.IdentExprKind:
		ident := expr.(*ast.IdentExpr).Value
		item, exists := env.Get(ident)
		if !exists {
			return std.ThrowException(env, fmt.Sprintf("variable '%s' does not exist", ident))
		}
		return item, state.Ok
	}
	panic("unexpected error")
}
