package utils

type Tuple[A any, B any] struct {
	First  A
	Second B
}

func TupleArrayToMap[A comparable, B any](arr []Tuple[A, B]) map[A]B {
	m := map[A]B{}

	for _, t := range arr {
		m[t.First] = t.Second
	}

	return m
}
