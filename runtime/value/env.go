package value

import (
	"bscript/parser/ast"
	"bscript/runtime/state"
)

type EvalStmt func(env *Env, stmt ast.Stmt) (Value, state.State)
type EvalFile func(file string) (Value, state.State)
type EvalSource func(source string) (Value, state.State)

type Variable struct {
	Value Value
	Const bool
}
type VariableMap map[string]*Variable

type GlobalContext struct {
	WorkingDir         string
	SourceFile, Source string
	Node               ast.Node
	Callstack          []string
	EvalStmt           EvalStmt
	EvalFile           EvalFile
	EvalSource         EvalSource
}

type Env struct {
	Global    *GlobalContext
	parent    *Env
	variables VariableMap
}

func NewEnv() *Env {
	return &Env{&GlobalContext{"", "", "", nil, make([]string, 0), nil, nil, nil}, nil, make(VariableMap)}
}
func (e *Env) NewChild() *Env {
	return &Env{e.Global, e, make(VariableMap)}
}
func (e *Env) SetCurrentNode(n ast.Node) {
	e.Global.Node = n
}
func (e *Env) GetCurrentNode() ast.Node {
	return e.Global.Node
}
func (e *Env) GetVariables() VariableMap {
	return e.variables
}
func (e *Env) Get(name string) (Value, bool) {
	item, exists := e.variables[name]
	if !exists {
		if e.parent != nil {
			return e.parent.Get(name)
		}
		return nil, false
	}
	return item.Value, true
}
func (e *Env) setOnly(name string, value Value) bool {
	item, exists := e.variables[name]
	if !exists {
		if e.parent != nil {
			return e.parent.setOnly(name, value)
		}
		return false
	}
	if item.Const {
		return false
	}
	item.Value = value
	return true
}
func (e *Env) Set(name string, value Value) bool {
	item, exists := e.variables[name]
	if !exists {
		if e.parent != nil {
			if e.parent.setOnly(name, value) {
				return true
			}
		}
		e.variables[name] = &Variable{value, false}
		return true
	}
	if item.Const {
		return false
	}
	item.Value = value
	return true
}
func (e *Env) Const(name string, value Value) bool {
	_, exists := e.variables[name]
	if exists {
		return false
	}
	e.variables[name] = &Variable{value, true}
	return true
}
