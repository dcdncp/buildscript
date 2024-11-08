package parser

import (
	"bscript/errors"
	"bscript/lexer"
	"bscript/parser/ast"
)

func ParseFile(file string) (*ast.ProgramStmt, []errors.Error, string) {
	tokens, _, source := lexer.TokenizeFile(file)
	p := NewParser(tokens, file, source)
	s, _ := p.ParseProgramStmt()
	return s, p.errors, source
}
func ParseSource(source string) (*ast.ProgramStmt, []errors.Error) {
	tokens, _ := lexer.TokenizeSource(source)
	p := NewParser(tokens, "<virtual>", source)
	s, _ := p.ParseProgramStmt()
	return s, p.errors
}
