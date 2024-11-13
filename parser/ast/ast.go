package ast

import "bscript/position"

type Node interface{
	position.Boundary
	String() string
}

type ExprKind int

const (
	IdentExprKind ExprKind = iota
	StringExprKind
	IntExprKind
	FloatExprKind
	BoolExprKind
	CallExprKind
	MemberExprKind
	IndexExprKind
	UnaryExprKind
	BinaryExprKind
	AssignExprKind
	TupleExprKind
	ArrayExprKind
	DefineExprKind
	ErrorExprKind
)

type Expr interface {
	Node
	Kind() ExprKind
}

type StmtKind int

const (
	ProgramStmtKind StmtKind = iota
	ConstStmtKind
	BlockStmtKind
	DefineStmtKind
	IfStmtKind
	ForStmtKind
	WhileStmtKind
	ExprStmtKind
	ReturnStmtKind
	BreakStmtKind
	ContinueStmtKind
	ThrowStmtKind
	ClassStmtKind
)

type Stmt interface {
	Node
	Kind() StmtKind
}
