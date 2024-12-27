package pkg

func Map[S ~[]T, T any, R any](s S, fn func(T) R) []R {
	if len(s) == 0 {
		return nil
	}

	result := make([]R, len(s))
	for i := range s {
		result[i] = fn(s[i])
	}
	return result
}
