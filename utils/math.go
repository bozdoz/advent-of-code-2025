package utils

// not sure what E type should be
func Sum[E ~int](slice []E) (out E) {
	for _, v := range slice {
		out += v
	}

	return
}

func SumFunc[T any, E ~int](slice []T, fun func(_ T) E) (out E) {
	for _, v := range slice {
		out += fun(v)
	}

	return
}
