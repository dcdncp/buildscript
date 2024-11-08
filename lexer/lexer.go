package lexer

import (
	"bscript/errors"
	"bscript/lexer/token"
	"bscript/position"
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	file               string
	source             string
	buffer             []rune
	errors             []errors.Error
	lastRow, lastCol   int
	position, row, col int
}

func NewLexer(file string, source string) *Lexer {
	return &Lexer{file, source, make([]rune, 0), make([]errors.Error, 0), 0, 0, 0, 0, 0}
}
func (l *Lexer) Source() string {
	return l.source
}
func (l *Lexer) SourceFile() string {
	return l.file
}
func (l *Lexer) LexerError(msg string) *errors.LexerError {
	e := errors.NewLexerError(l.file, l.source, msg,
		position.Position{Col: l.lastCol, Row: l.lastRow},
		position.Position{Col: l.col, Row: l.row})
	l.errors = append(l.errors, e)
	return e
}
func (l *Lexer) IsEOF() bool {
	return l.position >= len(l.source)
}
func (l *Lexer) It() rune {
	if l.IsEOF() {
		return '\000'
	}
	return rune(l.source[l.position])
}
func (l *Lexer) Next() rune {
	if l.position+1 > len(l.source) {
		return '\000'
	}
	return rune(l.source[l.position+1])
}
func (l *Lexer) Skip() rune {
	if l.IsEOF() {
		return '\000'
	}
	it := l.It()
	if it == '\n' {
		l.row += 1
		l.col = 0
	} else {
		l.col += 1
	}
	l.position += 1
	return it
}
func (l *Lexer) Eat() rune {
	if l.IsEOF() {
		return '\000'
	}
	it := l.Skip()
	l.buffer = append(l.buffer, it)
	return it
}
func (l *Lexer) NextToken() (token.Token, errors.Error) {
	l.lastCol = l.col
	l.lastRow = l.row
	lastPos := l.position
	kind, err := l.ParseToken()
	row := l.row
	col := l.col

	if err != nil {
		return token.Token{}, err
	}
	isLast := false
	for unicode.IsSpace(l.It()) && l.It() != '\n' {
		l.Skip()
	}
	if l.It() == '\n' {
		isLast = true
		l.Skip()
	}
	length := l.position - lastPos
	value := string(l.buffer)
	l.buffer = make([]rune, 0)
	return token.Token{
		Boundaries: position.CreateBoundaries(
			position.Position{Col: l.lastCol, Row: l.lastRow},
			position.Position{Col: col, Row: row}),
		Kind:       kind,
		Value:      value,
		Length:     length,
		IsLast:     isLast,
	}, nil
}
func (l *Lexer) ParseSpace() {
	for unicode.IsSpace(l.It()) {
		l.Skip()
	}
}
func (l *Lexer) ParseComment() {
	for l.It() != '\n' {
		l.Skip()
	}
}
func (l *Lexer) ParseIdent() {
	for unicode.IsLetter(l.It()) || unicode.IsDigit(l.It()) || l.It() == '_' {
		l.Eat()
	}
}
func (l *Lexer) ParseString() {
	for l.It() != '"' && l.It() != '\n' {
		it := l.Skip()
		if it == '\\' {
			it = l.Skip()
			if it == 'n' {
				it = '\n'
			} else if it == 't' {
				it = '\t'
			} else if it == '0' {
				it = '\000'
			}
		}
		l.buffer = append(l.buffer, it)
	}
}
func (l *Lexer) ParseInt() {
	for unicode.IsDigit(l.It()) {
		l.Eat()
	}
}
func (l *Lexer) CheckSymbol() bool {
	s := string(append(l.buffer, l.It()))
	for k := range token.SymbolMap {
		if strings.HasPrefix(k, s) {
			return true
		}
	}
	return false
}
func (l *Lexer) ParseToken() (token.Kind, errors.Error) {
	if l.IsEOF() {
		return token.EOF, nil
	}
	it := l.It()
	if unicode.IsSpace(it) {
		l.ParseSpace()
		return token.Space, nil
	} else if it == '/' && l.Next() == '/' {
		l.ParseComment()
		return token.Comment, nil
	} else if it == '"' {
		l.Skip()
		l.ParseString()
		if l.Skip() != '"' {
			return 0, l.LexerError("incorrect string literal")
		}
		return token.String, nil
	} else if unicode.IsLetter(it) || it == '_' {
		l.ParseIdent()
		ident := string(l.buffer)
		if ident == "true" || ident == "false" {
			return token.Bool, nil
		}
		k, exists := token.KeywordMap[ident]
		if exists {
			return k, nil
		}
		return token.Ident, nil
	} else if unicode.IsDigit(it) {
		l.ParseInt()
		if l.It() == rune(token.Dot) && unicode.IsDigit(l.Next()) {
			l.Eat()
			l.ParseInt()
			return token.Float, nil
		}
		return token.Int, nil
	} else if l.CheckSymbol() {
		l.Eat()
		for l.CheckSymbol() {
			l.Eat()
		}
		symbol := string(l.buffer)
		k, exists := token.SymbolMap[symbol]
		if exists {
			return k, nil
		}
		return 0, l.LexerError(fmt.Sprintf("unexpected '%s'", symbol))
	}
	return 0, l.LexerError(fmt.Sprintf("unexpected '%c'", l.Eat()))
}
