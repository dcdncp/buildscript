package tools

func Map[I any, O any](s []I, f func(I) O) []O {
	n := make([]O, 0)
	for _, e := range s {
		n = append(n, f(e))
	}
	return n
}

func AppendFront[I any](s []I, i I) []I {
	n := []I{i}
	return append(n, s...)
}
