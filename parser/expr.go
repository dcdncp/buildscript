package parser

import (
	"bscript/lexer/token"
	"bscript/parser/ast"
	"fmt"
	"slices"
	"strconv"
)

var logicTokens = []token.Kind{token.And, token.Or}
var cmpTokens = []token.Kind{token.Equals, token.NotEquals, token.Greater, token.Lower, token.GreaterEquals, token.LowerEquals}
var mulTokens = []token.Kind{token.Star, token.Slash}
var addTokens = []token.Kind{token.Plus, token.Minus}
var unaryTokens = []token.Kind{token.Not, token.Minus, token.Plus, token.TripleDot}
var mciTokens = []token.Kind{token.Dot, token.OpenParen, token.OpenBracket}

func (p *Parser) ParseExpr() (ast.Expr, bool) {
	return p.ParseAssignExpr()
}
func (p *Parser) ParseAssignExpr() (ast.Expr, bool) {
	lhs, err := p.ParseLogicExpr()
	for p.Match(token.Equal) {
		p.Eat()
		rhs, e := p.ParseLogicExpr()
		if !err {
			err = e
		}
		lhs = boundaries(ast.NewAssignExpr(lhs, rhs), lhs.Start(), rhs.End())
	}
	return lhs, err
}
func (p *Parser) ParseLogicExpr() (ast.Expr, bool) {
	lhs, err := p.ParseCmpExpr()
	for slices.Contains(logicTokens, p.It().Kind) {
		operator := p.Eat()
		rhs, e := p.ParseCmpExpr()
		if !err {
			err = e
		}
		lhs = boundaries(ast.NewBinaryExpr(operator, lhs, rhs), lhs.Start(), rhs.End())
	}
	return lhs, err
}
func (p *Parser) ParseCmpExpr() (ast.Expr, bool) {
	lhs, err := p.ParseMulExpr()
	for slices.Contains(cmpTokens, p.It().Kind) {
		operator := p.Eat()
		rhs, e := p.ParseMulExpr()
		if !err {
			err = e
		}
		lhs = boundaries(ast.NewBinaryExpr(operator, lhs, rhs), lhs.Start(), rhs.End())
	}
	return lhs, err
}
func (p *Parser) ParseMulExpr() (ast.Expr, bool) {
	lhs, err := p.ParseAddExpr()
	for slices.Contains(mulTokens, p.It().Kind) {
		operator := p.Eat()
		rhs, e := p.ParseAddExpr()
		if !err {
			err = e
		}
		lhs = boundaries(ast.NewBinaryExpr(operator, lhs, rhs), lhs.Start(), rhs.End())
	}
	return lhs, err
}
func (p *Parser) ParseAddExpr() (ast.Expr, bool) {
	lhs, err := p.ParseUnaryExpr()
	for slices.Contains(addTokens, p.It().Kind) {
		operator := p.Eat()
		rhs, e := p.ParseUnaryExpr()
		if !err {
			err = e
		}
		lhs = boundaries(ast.NewBinaryExpr(operator, lhs, rhs), lhs.Start(), rhs.End())
	}
	return lhs, err
}
func (p *Parser) ParseUnaryExpr() (ast.Expr, bool) {
	if slices.Contains(unaryTokens, p.It().Kind) {
		operator := p.Eat()
		operand, err := p.ParseUnaryExpr()
		e := boundaries(ast.NewUnaryExpr(operator, operand), operand.Start(), operand.End())
		return e, err
	}
	return p.ParseMCIExpr()
}
func (p *Parser) ParseMCIExpr() (ast.Expr, bool) {
	lhs, err := p.ParsePrimaryExpr()
	for slices.Contains(mciTokens, p.It().Kind) {
		start := lhs.Start()
		it := p.Eat()
		if it.Kind == token.Dot {
			field, e := p.Expect(token.Ident)
			if !err {
				err = e
			}
			lhs = ast.NewMemberExpr(lhs, field)
			lhs.SetEnd(field.End())
		} else if it.Kind == token.OpenParen {
			args := []ast.Expr{}
			for !p.Match(token.CloseParen) {
				expr, e := p.ParseExpr()
				if !err {
					err = e
				}
				args = append(args, expr)
				if p.Match(token.Comma) {
					p.Eat()
				} else {
					break
				}
			}
			end, e := p.Expect(token.CloseParen)
			if !err {
				err = e
			}
			lhs = ast.NewCallExpr(lhs, args)
			lhs.SetEnd(end.End())
		} else {
			index, e := p.ParseExpr()
			if !err {
				err = e
			}
			end, e := p.Expect(token.CloseBracket)
			if !err {
				err = e
			}
			lhs = ast.NewIndexExpr(lhs, index)
			lhs.SetEnd(end.End())
		}
		lhs.SetStart(start)
	}
	return lhs, err
}
func (p *Parser) ParsePrimaryExpr() (ast.Expr, bool) {
	if p.Match(token.Int) {
		it := p.Eat()
		val, _ := strconv.ParseInt(it.Value, 10, 64)
		return boundaries(ast.NewIntExpr(val), it.Start(), it.End()), false
	} else if p.Match(token.Float) {
		it := p.Eat()
		val, _ := strconv.ParseFloat(it.Value, 64)
		return boundaries(ast.NewFloatExpr(val), it.Start(), it.End()), false
	} else if p.Match(token.Bool) {
		it := p.Eat()
		var e ast.Expr
		if it.Value == "true" {
			e = ast.NewBoolExpr(true)
		} else {
			e = ast.NewBoolExpr(false)
		}
		return boundaries(e, it.Start(), it.End()), false
	} else if p.Match(token.String) {
		it := p.Eat()
		e := ast.NewStringExpr(it.Value)
		return boundaries(e, it.Start(), it.End()), false
	} else if p.Match(token.Ident) {
		it := p.Eat()
		e := ast.NewIdentExpr(it.Value)
		return boundaries(e, it.Start(), it.End()), false
	} else if p.Match(token.OpenParen) {
		start := p.Eat()
		expr, err := p.ParseExpr()
		if err {
			return expr, err
		}
		if p.Match(token.CloseParen) {
			p.Eat()
			return expr, false
		}
		values := []ast.Expr{expr}
		for p.Match(token.Comma) {
			p.Eat()
			expr, e := p.ParseExpr()
			if !err {
				err = e
			}
			values = append(values, expr)
		}
		end, e := p.Expect(token.CloseParen)
		if !err {
			err = e
		}
		return boundaries(ast.NewTupleExpr(values), start.Start(), end.End()), err
	} else if p.Match(token.OpenBracket) {
		err := false
		start := p.Eat()
		values := []ast.Expr{}
		for !p.Match(token.CloseBracket) {
			expr, e := p.ParseExpr()
			if !err {
				err = e
			}
			values = append(values, expr)
			if p.Match(token.Comma) {
				p.Eat()
			} else {
				break
			}
		}
		end, e := p.Expect(token.CloseBracket)
		if !err {
			err = e
		}
		return boundaries(ast.NewArrayExpr(values), start.Start(), end.End()), err
	} else if p.Match(token.Define) {
		start := p.Eat()
		err := false
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
			ast.NewDefineExpr(params, variadic, body), start.Start(), body.End()), err
	
	}
	it := p.Eat()
	p.ParserError(fmt.Sprintf("unexpected '%s'", it.Value), it.Start(), it.End())
	return boundaries(ast.NewErrorExpr(), it.Start(), it.End()), true
}
