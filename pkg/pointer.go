package pkg

func Pointer[T any](v T) *T {
	return &v
}
