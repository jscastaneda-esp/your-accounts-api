package slices

func Find[T any](s []T, fn func(e T) bool) (e T) {
	for _, el := range s {
		if fn(el) {
			e = el
			return
		}
	}

	return
}
