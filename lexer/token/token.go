package token

import (
	"bscript/position"
	"fmt"
)

type Kind int

const (
	// Literals
	Ident Kind = iota
	String
	Int
	Float
	Bool

	// Keywords
	Const
	Throw
	For
	If
	Else
	Define
	End
	While
	And
	Or
	Not
	In
	Return
	Break
	Continue

	// Operators
	Equal
	Equals
	NotEquals
	Greater
	Lower
	GreaterEquals
	LowerEquals
	Arrow
	Minus
	Plus
	Star
	Slash
	Dot
	TripleDot
	Colon
	Comma
	OpenParen
	CloseParen
	OpenBracket
	CloseBracket
	OpenBrace
	CloseBrace

	// Other
	Comment
	Space
	EOF
)

func (k Kind) String() string {
	switch k {
	case Ident:
		return "ident"
	case Int:
		return "int"
	case String:
		return "string"
	case Float:
		return "float"
	case Bool:
		return "bool"

	case Throw:
		return "throw"
	case Const:
		return "const"
	case For:
		return "for"
	case If:
		return "if"
	case Else:
		return "else"
	case Define:
		return "define"
	case End:
		return "end"
	case While:
		return "while"
	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"
	case In:
		return "in"
	case Return:
		return "return"
	case Break:
		return "break"
	case Continue:
		return "continue"

	case Equal:
		return "="
	case Equals:
		return "=="
	case NotEquals:
		return "!="
	case Greater:
		return ">"
	case Lower:
		return "<"
	case GreaterEquals:
		return ">="
	case LowerEquals:
		return "<="
	case Arrow:
		return "<-"
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Star:
		return "*"
	case Slash:
		return "/"
	case Dot:
		return "."
	case TripleDot:
		return "..."
	case Colon:
		return ":"
	case Comma:
		return ","
	case OpenParen:
		return "("
	case CloseParen:
		return ")"
	case OpenBracket:
		return "["
	case CloseBracket:
		return "]"
	case OpenBrace:
		return "{"
	case CloseBrace:
		return "}"

	case Comment:
		return "comment"
	case Space:
		return "space"
	case EOF:
		return "eof"
	}
	return "unreachable"
}

var SymbolMap = map[string]Kind{
	"=":   Equal,
	"+":   Plus,
	"-":   Minus,
	"*":   Star,
	"/":   Slash,
	"(":   OpenParen,
	")":   CloseParen,
	"{":   OpenBrace,
	"}":   CloseBrace,
	"[":   OpenBracket,
	"]":   CloseBracket,
	".":   Dot,
	":":   Colon,
	"<-":  Arrow,
	"...": TripleDot,
	",":   Comma,
	"==":  Equals,
	"!=":  NotEquals,
	">":   Greater,
	"<":   Lower,
	">=":  GreaterEquals,
	"<=":  LowerEquals,
}

var KeywordMap = map[string]Kind{
	"const":    Const,
	"throw":    Throw,
	"for":      For,
	"define":   Define,
	"if":       If,
	"else":     Else,
	"end":      End,
	"while":    While,
	"not":      Not,
	"or":       Or,
	"and":      And,
	"in":       In,
	"continue": Continue,
	"break":    Break,
	"return":   Return,
}

type Token struct {
	position.Boundaries
	Kind       Kind
	Value      string
	Length     int
	IsLast     bool
}

func (t Token) String() string {
	if len(t.Value) > 0 {
		if t.Kind < 5 {
			return fmt.Sprintf("Token{%s, '%s'}", t.Kind, t.Value)
		} else {
			return fmt.Sprintf("Token{'%s'}", t.Value)
		}
	} else {
		return fmt.Sprintf("Token{%s}", t.Kind)
	}
}
