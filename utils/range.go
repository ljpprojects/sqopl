package utils

type Range[T any] struct {
	Start T
	End   T
}

func InitRange[T any](start T, end T) Range[T] {
	return Range[T]{
		Start: start,
		End:   end,
	}
}
