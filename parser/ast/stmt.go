package ast

import (
	"bscript/lexer/token"
	"bscript/position"
	"bscript/tools"
	"fmt"
	"strings"
)

type ProgramStmt struct {
	position.Boundaries
	Body []Stmt
}

func NewProgramStmt(body []Stmt) *ProgramStmt {
	return &ProgramStmt{Body: body}
}
func (ProgramStmt) Kind() StmtKind {
	return ProgramStmtKind
}
func (ProgramStmt) String() string {
	return "<program>"
}

type ConstStmt struct {
	position.Boundaries
	Id    Expr
	Value Expr
}

func NewConstStmt(id Expr, value Expr) *ConstStmt {
	return &ConstStmt{Id: id, Value: value}
}
func (ConstStmt) Kind() StmtKind {
	return ConstStmtKind
}
func (s *ConstStmt) String() string {
	return fmt.Sprintf("<const %s = %s>", s.Id.String(), s.Value.String())
}

type BlockStmt struct {
	position.Boundaries
	Body []Stmt
}

func NewBlockStmt(body []Stmt) *BlockStmt {
	return &BlockStmt{Body: body}
}
func (BlockStmt) Kind() StmtKind {
	return BlockStmtKind
}
func (BlockStmt) String() string {
	return "<body>"
}

type DefineStmt struct {
	position.Boundaries
	Ident    token.Token
	Params   []token.Token
	Variadic bool
	Body     Stmt
}

func NewDefineStmt(ident token.Token, params []token.Token, variadic bool, body Stmt) *DefineStmt {
	return &DefineStmt{Ident: ident, Params: params, Variadic: variadic, Body: body}
}
func (DefineStmt) Kind() StmtKind {
	return DefineStmtKind
}
func (s *DefineStmt) String() string {
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
	return fmt.Sprintf("<define %s(%s)>", s.Ident.Value, str)
}

type IfStmt struct {
	position.Boundaries
	Condition Expr
	Primary   Stmt
	Secondary Stmt
}

func NewIfStmt(condition Expr, primary Stmt, secondary Stmt) *IfStmt {
	return &IfStmt{Condition: condition, Primary: primary, Secondary: secondary}
}
func (IfStmt) Kind() StmtKind {
	return IfStmtKind
}
func (s *IfStmt) String() string {
	return fmt.Sprintf("<if %s>", s.Condition)
}

type WhileStmt struct {
	position.Boundaries
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{Condition: condition, Body: body}
}
func (WhileStmt) Kind() StmtKind {
	return WhileStmtKind
}
func (s *WhileStmt) String() string {
	return fmt.Sprintf("<while %s>", s.Condition)
}

type ForStmt struct {
	position.Boundaries
	Ident token.Token
	In    Expr
	Body  Stmt
}

func NewForStmt(ident token.Token, in Expr, body Stmt) *ForStmt {
	return &ForStmt{Ident: ident, In: in, Body: body}
}
func (ForStmt) Kind() StmtKind {
	return ForStmtKind
}
func (s *ForStmt) String() string {
	return fmt.Sprintf("<for %s in %s>", s.Ident.Value, s.In)
}

type ExprStmt struct {
	position.Boundaries
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{Expr: expr}
}
func (ExprStmt) Kind() StmtKind {
	return ExprStmtKind
}
func (s *ExprStmt) String() string {
	return s.Expr.String()
}

type BreakStmt struct {
	position.Boundaries
}

func NewBreakStmt() *BreakStmt {
	return &BreakStmt{}
}
func (BreakStmt) Kind() StmtKind {
	return BreakStmtKind
}
func (BreakStmt) String() string {
	return "<break>"
}

type ContinueStmt struct {
	position.Boundaries
}

func NewContinueStmt() *ContinueStmt {
	return &ContinueStmt{}
}
func (ContinueStmt) Kind() StmtKind {
	return ContinueStmtKind
}
func (ContinueStmt) String() string {
	return "<continue>"
}

type ReturnStmt struct {
	position.Boundaries
	Value Expr
}

func NewReturnStmt(value Expr) *ReturnStmt {
	return &ReturnStmt{Value: value}
}
func (ReturnStmt) Kind() StmtKind {
	return ReturnStmtKind
}
func (s *ReturnStmt) String() string {
	return fmt.Sprintf("<return %s>", s.Value)
}

type ThrowStmt struct {
	position.Boundaries
	Value Expr
}

func NewThrowStmt(value Expr) *ThrowStmt {
	return &ThrowStmt{Value: value}
}
func (ThrowStmt) Kind() StmtKind {
	return ThrowStmtKind
}
func (s *ThrowStmt) String() string {
	return fmt.Sprintf("<throw %s>", s.Value)
}

type ClassStmt struct {
	position.Boundaries
	Name token.Token
	Body Stmt
}

func NewClassStmt(name token.Token, body Stmt) *ClassStmt {
	return &ClassStmt{Name: name, Body: body}
}
func (ClassStmt) Kind() StmtKind {
	return ClassStmtKind
}
func (s *ClassStmt) String() string {
	return fmt.Sprintf("<class %s>", s.Name)
}
