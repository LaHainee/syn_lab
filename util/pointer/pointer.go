package pointer

func Val[T comparable](target *T) T {
	var result T
	if target == nil {
		return result
	}

	return *target
}

func To[T any](target T) *T {
	return &target
}
