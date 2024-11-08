package ast

import (
	"bscript/lexer/token"
	"bscript/position"
	"bscript/tools"
	"fmt"
	"strings"
)

type IntExpr struct {
	position.Boundaries
	Value int64
}

func NewIntExpr(value int64) *IntExpr {
	return &IntExpr{Value: value}
}
func (IntExpr) Kind() ExprKind {
	return IntExprKind
}
func (e *IntExpr) String() string {
	return fmt.Sprintf("%d", e.Value)
}

type FloatExpr struct {
	position.Boundaries
	Value float64
}

func NewFloatExpr(value float64) *FloatExpr {
	return &FloatExpr{Value: value}
}
func (FloatExpr) Kind() ExprKind {
	return FloatExprKind
}
func (e *FloatExpr) String() string {
	return fmt.Sprintf("%f", e.Value)
}

type StringExpr struct {
	position.Boundaries
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{Value: value}
}
func (StringExpr) Kind() ExprKind {
	return StringExprKind
}
func (e *StringExpr) String() string {
	return fmt.Sprintf("\"%s\"", e.Value)
}

type BoolExpr struct {
	position.Boundaries
	Value bool
}

func NewBoolExpr(value bool) *BoolExpr {
	return &BoolExpr{Value: value}
}
func (BoolExpr) Kind() ExprKind {
	return BoolExprKind
}
func (e *BoolExpr) String() string {
	if e.Value {
		return "true"
	} else {
		return "false"
	}
}

type IdentExpr struct {
	position.Boundaries
	Value string
}

func NewIdentExpr(value string) *IdentExpr {
	return &IdentExpr{Value: value}
}
func (IdentExpr) Kind() ExprKind {
	return IdentExprKind
}
func (e *IdentExpr) String() string {
	return e.Value
}

type UnaryExpr struct {
	position.Boundaries
	Operator token.Token
	Operand  Expr
}

func NewUnaryExpr(operator token.Token, operand Expr) *UnaryExpr {
	return &UnaryExpr{Operator: operator, Operand: operand}
}
func (UnaryExpr) Kind() ExprKind {
	return UnaryExprKind
}
func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", e.Operator.Value, e.Operand)
}

type BinaryExpr struct {
	position.Boundaries
	Operator token.Token
	LHS, RHS Expr
}

func NewBinaryExpr(operator token.Token, lhs, rhs Expr) *BinaryExpr {
	return &BinaryExpr{Operator: operator, LHS: lhs, RHS: rhs}
}
func (BinaryExpr) Kind() ExprKind {
	return BinaryExprKind
}
func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.LHS, e.Operator.Value, e.RHS)
}

type AssignExpr struct {
	position.Boundaries
	LHS, RHS Expr
}

func NewAssignExpr(lhs, rhs Expr) *AssignExpr {
	return &AssignExpr{LHS: lhs, RHS: rhs}
}
func (AssignExpr) Kind() ExprKind {
	return AssignExprKind
}
func (e *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", e.LHS, e.RHS)
}

type CallExpr struct {
	position.Boundaries
	Target Expr
	Args   []Expr
}

func NewCallExpr(target Expr, args []Expr) *CallExpr {
	return &CallExpr{Target: target, Args: args}
}
func (CallExpr) Kind() ExprKind {
	return CallExprKind
}
func (e *CallExpr) String() string {
	return fmt.Sprintf(
		"%s(%s)", e.Target, strings.Join(tools.Map(e.Args, Expr.String), ", "))
}

type MemberExpr struct {
	position.Boundaries
	Target Expr
	Field  token.Token
}

func NewMemberExpr(target Expr, field token.Token) *MemberExpr {
	return &MemberExpr{Target: target, Field: field}
}
func (MemberExpr) Kind() ExprKind {
	return MemberExprKind
}
func (e *MemberExpr) String() string {
	return fmt.Sprintf("%s.%s", e.Target, e.Field.Value)
}

type IndexExpr struct {
	position.Boundaries
	Target Expr
	Index  Expr
}

func NewIndexExpr(target Expr, index Expr) *IndexExpr {
	return &IndexExpr{Target: target, Index: index}
}
func (IndexExpr) Kind() ExprKind {
	return IndexExprKind
}
func (e *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", e.Target, e.Index)
}

type TupleExpr struct {
	position.Boundaries
	Values []Expr
}

func NewTupleExpr(values []Expr) *TupleExpr {
	return &TupleExpr{Values: values}
}
func (TupleExpr) Kind() ExprKind {
	return TupleExprKind
}
func (e *TupleExpr) String() string {
	return fmt.Sprintf(
		"(%s)", strings.Join(tools.Map(e.Values, Expr.String), ", "))
}

type ArrayExpr struct {
	position.Boundaries
	Values []Expr
}

func NewArrayExpr(values []Expr) *ArrayExpr {
	return &ArrayExpr{Values: values}
}
func (ArrayExpr) Kind() ExprKind {
	return ArrayExprKind
}
func (e *ArrayExpr) String() string {
	return fmt.Sprintf(
		"[%s]", strings.Join(tools.Map(e.Values, Expr.String), ", "))
}

type ErrorExpr struct {
	position.Boundaries
}

type DefineExpr struct {
	position.Boundaries
	Params   []token.Token
	Variadic bool
	Body     Stmt
}

func NewDefineExpr(params []token.Token, variadic bool, body Stmt) *DefineExpr {
	return &DefineExpr{Params: params, Variadic: variadic, Body: body}
}
func (DefineExpr) Kind() ExprKind {
	return DefineExprKind
}
func (s *DefineExpr) String() string {
	var str string
	if s.Variadic {
		str = strings.Join(
			tools.Map(s.Params[:len(s.Params)-2],
				func(t token.Token) string { return t.Value }), ", ")
		str += ", ..." + s.Params[len(s.Params)-1].Value
	} else {
		str = strings.Join(
			tools.Map(s.Params,
				func(t token.Token) string { return t.Value }), ", ")
	}
	return fmt.Sprintf("<define (%s)>", str)
}

func NewErrorExpr() *ErrorExpr {
	return &ErrorExpr{}
}
func (ErrorExpr) Kind() ExprKind {
	return ErrorExprKind
}
func (ErrorExpr) String() string {
	return "<error>"
}