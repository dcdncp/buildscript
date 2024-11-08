package state

type State int

const (
	Ok State = iota
	Error
	Return
	Break
	Continue
)

func (s State) IsOkay() bool {
	return s == Ok
}
func (s State) IsNotOkay() bool {
	return s != Ok
}

func (s State) String() string {
	switch s {
	case Ok:
		return "<state Ok>"
	case Error:
		return "<state Error>"
	case Return:
		return "<state Return>"
	case Break:
		return "<state Break>"
	case Continue:
		return "<state Continue>"
	default:
		return "unreachable"
	}
}
