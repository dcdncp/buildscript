package position

type Position struct {
	Col, Row int
}

type Boundary interface {
	Start() Position
	End() Position
	SetStart(p Position)
	SetEnd(p Position)
}

type Boundaries struct {
	start, end Position
}

func CreateBoundaries(start, end Position) Boundaries {
	return Boundaries{start, end}
}


func (b *Boundaries) Start() Position {
	return b.start
}
func (b *Boundaries) End() Position {
	return b.end
}
func (b *Boundaries) SetStart(p Position) {
	b.start = p
}
func (b *Boundaries) SetEnd(p Position) {
	b.end = p
}
