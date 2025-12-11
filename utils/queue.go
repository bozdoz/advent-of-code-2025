package utils

// gets last (LIFO)
func Pop[T any](s *[]T) *T {
	if len(*s) == 0 {
		return nil
	}

	end := len(*s) - 1
	val := (*s)[end]
	*s = (*s)[:end]

	return &val
}

// gets last (LIFO), without checking length
func MustPop[T any](s *[]T) T {
	end := len(*s) - 1
	val := (*s)[end]
	*s = (*s)[:end]

	return val
}

// gets first (FIFO)
func Shift[T any](s *[]T) *T {
	if len(*s) == 0 {
		return nil
	}

	val := (*s)[0]
	*s = (*s)[1:]

	return &val
}

// gets first (FIFO), without checking length
func MustShift[T any](s *[]T) T {
	val := (*s)[0]
	*s = (*s)[1:]

	return val
}
