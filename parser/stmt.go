package parser

import (
	"bscript/lexer/token"
	"bscript/parser/ast"
	"bscript/position"
)

func (p *Parser) ParseProgramStmt() (*ast.ProgramStmt, bool) {
	err := false
	body := make([]ast.Stmt, 0)
	start := p.It()
	for !p.IsEOF() {
		s, e := p.ParseStmt()
		if !err {
			err = e
		}
		body = append(body, s)
	}
	program := ast.NewProgramStmt(body)
	end := p.It()
	return boundaries(program, start.Start(), end.End()), err
}
func (p *Parser) ParseStmt() (ast.Stmt, bool) {
	if p.Match(token.Define) {
		return p.ParseDefineStmt()
	} else if p.Match(token.Const) {
		return p.ParseConstStmt()
	} else if p.Match(token.If) {
		return p.ParseIfStmt()
	} else if p.Match(token.While) {
		return p.ParseWhileStmt()
	} else if p.Match(token.For) {
		return p.ParseForStmt()
	} else if p.Match(token.Continue) {
		it := p.Eat()
		return boundaries(ast.NewContinueStmt(), it.Start(), it.End()), false
	} else if p.Match(token.Break) {
		it := p.Eat()
		return boundaries(ast.NewBreakStmt(), it.Start(), it.End()), false
	} else if p.Match(token.Return) {
		it := p.Eat()
		end := it.End()
		var value ast.Expr
		err := false
		if !it.IsLast {
			value, err = p.ParseExpr()
			end = value.End()
		}
		return boundaries(ast.NewReturnStmt(value), it.Start(), end), err
	} else if p.Match(token.Throw) {
		it := p.Eat()
		expr, err := p.ParseExpr()
		return boundaries(ast.NewThrowStmt(expr), it.Start(), expr.End()), err
	} else {
		expr, err := p.ParseExpr()
		return boundaries(ast.NewExprStmt(expr), expr.Start(), expr.End()), err
	}
}
func (p *Parser) ParseBlockStmt() (*ast.BlockStmt, bool) {
	err := false
	start := p.It()
	body := []ast.Stmt{}
	for !p.Match(token.End) && !p.IsEOF() {
		s, e := p.ParseStmt()
		if !err {
			err = e
		}
		body = append(body, s)
	}
	end, e := p.Expect(token.End)
	if !err {
		err = e
	}
	return boundaries(ast.NewBlockStmt(body), start.Start(), end.End()), err
}
func (p *Parser) ParseDefineStmt() (*ast.DefineStmt, bool) {
	start := p.Eat()
	ident, err := p.Expect(token.Ident)
	params := make([]token.Token, 0)
	variadic := false
	p.Expect(token.OpenParen)
	for !p.Match(token.CloseParen) {
		if p.Match(token.TripleDot) {
			p.Eat()
			variadic = true
		}
		param, e := p.Expect(token.Ident)
		if !err {
			err = e
		}
		params = append(params, param)
		if variadic && !p.Match(token.CloseParen) {
			p.ParserError("variadic parameter must be a last one", param.Start(), param.End())
		} else if p.Match(token.Comma) {
			p.Eat()
		} else {
			break
		}
	}
	p.Expect(token.CloseParen)
	body, e := p.ParseBlockStmt()
	if !err {
		err = e
	}
	return boundaries(
		ast.NewDefineStmt(ident, params, variadic, body), start.Start(), body.End()), err
}
func (p *Parser) ParseConstStmt() (*ast.ConstStmt, bool) {
	start := p.Eat()
	id, err := p.Expect(token.Ident)
	p.Expect(token.Equal)
	value, e := p.ParseExpr()
	if !err {
		err = e
	}
	return boundaries(
		ast.NewConstStmt(ast.NewIdentExpr(id.Value), value), start.Start(), value.End()), err
}
func (p *Parser) ParseIfStmt() (*ast.IfStmt, bool) {
	start := p.Eat()
	condition, err := p.ParseExpr()
	body := make([]ast.Stmt, 0)
	primaryStart := p.It()
	for !p.Match(token.End) && !p.Match(token.Else) && !p.IsEOF() {
		s, e := p.ParseStmt()
		if !err {
			err = e
		}
		body = append(body, s)
	}
	var end position.Position
	var stmt *ast.IfStmt
	if p.Match(token.End) {
		it := p.Eat()
		end = it.End()
		primary := ast.NewBlockStmt(body)
		primary.SetStart(primaryStart.Start())
		primary.SetEnd(end)
		stmt = ast.NewIfStmt(condition, primary, nil)
	} else if p.Match(token.Else) {
		it := p.Eat()
		end = it.End()
		primary := ast.NewBlockStmt(body)
		primary.SetStart(primaryStart.Start())
		primary.SetEnd(end)
		if p.Match(token.If) {
			secondary, e := p.ParseIfStmt()
			if !err {
				err = e
			}
			stmt = ast.NewIfStmt(condition, primary, secondary)
			end = secondary.End()
		} else {
			secondary, e := p.ParseBlockStmt()
			if !err {
				err = e
			}
			stmt = ast.NewIfStmt(condition, primary, secondary)
			end = secondary.End()
		}
	} else {
		it := p.It()
		p.ParserError("expected end or else", it.Start(), it.End())
	}
	stmt.SetStart(start.Start())
	stmt.SetEnd(end)
	return stmt, err
}
func (p *Parser) ParseForStmt() (*ast.ForStmt, bool) {
	start := p.Eat()
	ident, err := p.Expect(token.Ident)
	p.Expect(token.In)
	in, e := p.ParseExpr()
	if !err {
		err = e
	}
	body, e := p.ParseBlockStmt()
	if !err {
		err = e
	}
	return boundaries(ast.NewForStmt(ident, in, body), start.Start(), body.End()), err
}
func (p *Parser) ParseWhileStmt() (*ast.WhileStmt, bool) {
	start := p.Eat()
	condition, err := p.ParseExpr()
	body, e := p.ParseBlockStmt()
	if !err {
		err = e
	}
	return boundaries(ast.NewWhileStmt(condition, body), start.Start(), body.End()), err
}
