package gosmarty

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Ptr[T any](a T) *T {
	return &a
}
