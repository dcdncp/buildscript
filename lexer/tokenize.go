package lexer

import (
	"bscript/errors"
	"bscript/lexer/token"
	"os"
)

func TokenizeFile(path string) ([]token.Token, []errors.Error, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	l := NewLexer(path, string(data))
	tokens := make([]token.Token, 0)
	for !l.IsEOF() {
		tok, err := l.NextToken()
		if err != nil {
			continue
		}
		if tok.Kind == token.Space || tok.Kind == token.Comment {
			continue
		}
		tokens = append(tokens, tok)
	}
	tok, _ := l.NextToken()
	tokens = append(tokens, tok)
	return tokens, l.errors, l.source
}
func TokenizeSource(source string) ([]token.Token, []errors.Error) {
	l := NewLexer("<virtual>", source)
	tokens := make([]token.Token, 0)
	for !l.IsEOF() {
		tok, err := l.NextToken()
		if err != nil {
			continue
		}
		if tok.Kind == token.Space || tok.Kind == token.Comment {
			continue
		}
		tokens = append(tokens, tok)
	}
	tok, _ :=  l.NextToken()
	tokens = append(tokens, tok)
	return tokens, l.errors
}