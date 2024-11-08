package parser

import (
	"bscript/errors"
	"bscript/lexer/token"
	"bscript/position"
	"fmt"
)

func boundaries[T position.Boundary](n T, start, end position.Position) T {
	n.SetStart(start)
	n.SetEnd(end)
	return n
}

type Parser struct {
	file     string
	source   string
	tokens   []token.Token
	position int
	errors   []errors.Error
}

func NewParser(tokens []token.Token, file, source string) *Parser {
	return &Parser{file, source, tokens, 0, make([]errors.Error, 0)}
}

func (p *Parser) Source() string {
	return p.source
}
func (p *Parser) SkipLine() {
	for !p.It().IsLast && !p.IsEOF() {
		p.Eat()
	}
}
func (p *Parser) ParserError(msg string, start, end position.Position) bool {
	e := errors.NewParserError(p.file, p.source, msg, start, end)
	p.errors = append(p.errors, e)
	return false
}
func (p *Parser) IsEOF() bool {
	return p.position+1 >= len(p.tokens)
}
func (p *Parser) It() token.Token {
	return p.tokens[p.position]
}
func (p *Parser) Eat() token.Token {
	it := p.It()
	if !p.IsEOF() {
		p.position += 1
	}
	return it
}
func (p *Parser) Match(kind token.Kind) bool {
	return p.It().Kind == kind
}
func (p *Parser) Expect(kind token.Kind) (token.Token, bool) {
	it := p.It()
	if it.Kind != kind {
		p.ParserError(
			fmt.Sprintf("unexpected '%s', required %s", it.Value, kind), it.Start(), it.End())
		return it, true
	}
	return p.Eat(), false
}
