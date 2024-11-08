package errors

import (
	"bscript/errors/kind"
	"bscript/position"
	"fmt"
	"strings"
)

var defaultPad = 5
var delimeter = " | "

func errorPlaceString(err Error) string {
	out := ""
	lines := strings.Split(err.Source(), "\n")
	linesCount := err.End().Row - err.Start().Row + 1
	idx := err.Start().Row
	start, end := err.Start().Col, err.End().Col
	if idx > 0 {
		lineStr := fmt.Sprint(idx)
		lineStr = strings.Repeat(" ", defaultPad-len(lineStr)) + lineStr + delimeter + lines[idx-1]
		out += lineStr + "\n"
	}
	if linesCount == 1 {
		lineStr := fmt.Sprint(idx + 1)
		lineStr = strings.Repeat(" ", defaultPad-len(lineStr)) + lineStr + delimeter + lines[idx]
		out += lineStr + "\n"
	} else {
		for i := 0; i < linesCount; i++ {
			index := idx + i
			line := lines[index]
			lsLine := strings.TrimLeft(line, " ")
			rsLine := strings.TrimRight(line, " ")
			startC := len(line) - len(lsLine)
			if startC < start {
				start = startC
			}
			endC := len(rsLine)
			if endC > end {
				end = endC
			}
			lineStr := fmt.Sprint(index + 1)
			lineStr = strings.Repeat(" ", defaultPad-len(lineStr)) + lineStr + delimeter + line
			out += lineStr + "\n"
		}
	}
	length := end - start - 1
	if length < 0 {
		length = 0
	}
	hl := strings.Repeat(" ", defaultPad) + delimeter + strings.Repeat(" ", start) + "^" + strings.Repeat("~", length)
	out += hl + "\n"
	return out
}

func errorLabelString(err Error, prefix string) string {
	return fmt.Sprintf("%s:%d:%d: %s: %s\n", err.SourceFile(), err.Start().Row+1, err.Start().Col+1, prefix, err.Message())
}

type Error interface {
	position.Boundary
	Kind() kind.Kind
	Source() string
	SourceFile() string
	Message() string
}

type CommonError struct {
	position.Boundaries
	file, source, message string
}

func (b *CommonError) Source() string {
	return b.source
}
func (b *CommonError) SourceFile() string {
	return b.file
}
func (b *CommonError) Message() string {
	return b.message
}
func (b *CommonError) String() string {
	out := errorLabelString(b, "error")
	out += errorPlaceString(b)
	return out
}
func (b *CommonError) Kind() kind.Kind {
	return kind.Common
}

type LexerError struct {
	CommonError
}

func NewLexerError(file, source, message string, start, end position.Position) *LexerError {
	return &LexerError{
		CommonError{position.CreateBoundaries(start, end), file, source, message}}
}
func (e *LexerError) Kind() kind.Kind {
	return kind.Lexer
}
func (e *LexerError) String() string {
	out := errorLabelString(e, "lexer error")
	out += errorPlaceString(e)
	return out
}

type ParserError struct {
	CommonError
}

func NewParserError(file, source, message string, start, end position.Position) *ParserError {
	return &ParserError{CommonError{position.CreateBoundaries(start, end), file, source, message}}
}
func (e *ParserError) Kind() kind.Kind {
	return kind.Parser
}
func (e *ParserError) String() string {
	out := errorLabelString(e, "parser error")
	out += errorPlaceString(e)
	return out
}

type RuntimeError struct {
	CommonError
}

func NewRuntimeError(file, source, message string, start, end position.Position) *RuntimeError {
	return &RuntimeError{CommonError{position.CreateBoundaries(start, end), file, source, message}}
}
func (e *RuntimeError) Kind() kind.Kind {
	return kind.Runtime
}
func (e *RuntimeError) String() string {
	out := errorLabelString(e, "runtime error")
	out += errorPlaceString(e)
	return out
}
