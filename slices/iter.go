package slices

func Map[T, U any](slice []T, mapper func(item T, index int) U) []U {
	var result []U
	for i, v := range slice {
		result = append(result, mapper(v, i))
	}
	return result
}
