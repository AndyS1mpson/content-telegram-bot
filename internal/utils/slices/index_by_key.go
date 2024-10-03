package slices

// IndexByKey индексирует список по ключу, получаемому на основе этой структуры.
func IndexByKey[T any, K comparable](s []T, keyExtractor func(elem T) K) map[K]T {
	res := make(map[K]T, len(s))

	if keyExtractor == nil {
		return res
	}

	for _, elem := range s {
		res[keyExtractor(elem)] = elem
	}

	return res
}

// IndexByKeyValue индексирует список по ключу из структуры и модифицирует значение.
func IndexByKeyValue[T any, K comparable, V any](
	slice []T,
	keyFunc func(elem T) K,
	mapper func(elem T) V,
) map[K]V {
	res := make(map[K]V, len(slice))
	for _, elem := range slice {
		res[keyFunc(elem)] = mapper(elem)
	}

	return res
}
