package kind

type Kind int

const (
	Common Kind = iota
	Lexer 
	Parser
	Runtime
)