package slices

func Filter[T any](s []T, fn func(e T) bool) []T {
	e := []T{}

	for _, el := range s {
		if fn(el) {
			e = append(e, el)
		}
	}

	return e
}

func Find[T any](s []T, fn func(e T) bool) (e T) {
	for _, el := range s {
		if fn(el) {
			e = el
			return
		}
	}

	return
}
