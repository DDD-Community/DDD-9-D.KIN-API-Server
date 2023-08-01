package typex

func P[T any](v T) *T {
	return &v
}

func PrimitiveDeepEq[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	} else if a != nil && b != nil {
		return *a == *b
	}

	return false
}

type SliceComparable[T comparable] []T

func (list SliceComparable[T]) Contains(v T) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}

	return false
}

func Contains[T comparable](list []T, value T) bool {
	return SliceComparable[T](list).Contains(value)
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func SkipError[T any](t T, _ error) T {
	return t
}
