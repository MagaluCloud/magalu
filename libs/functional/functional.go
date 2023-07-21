package functional

func Transform[T any, U any](source []T, transform func(T) U) []U {
	result := make([]U, 0, len(source))
	for _, o := range source {
		result = append(result, transform(o))
	}
	return result
}

func TransformMap[T comparable, U any, V any](source map[T]U, transform func(key T, value U) V) []V {
	length := len(source)
	result := make([]V, length)

	i := 0
	for key, value := range source {
		result[i] = transform(key, value)
		i++
	}
	return result
}

func Reduce[T any, U any](source []T, into U, reduce func(o T, current *U)) U {
	for _, o := range source {
		reduce(o, &into)
	}
	return into
}

func Merge[T comparable, U any](base map[T]U, merger map[T]U) map[T]U {
	result := base

	for key, value := range merger {
		result[key] = value
	}

	return result
}

func Contains[T comparable](slice []T, element T) bool {
	for _, o := range slice {
		if o == element {
			return true
		}
	}

	return false
}
