package typex

func P[T any](v T) *T {
	return &v
}
