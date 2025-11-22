package util

func MapOfRefsToInterfaces[K comparable, T any, I any](value map[K]*T) map[K]I {
	newMap := make(map[K]I, len(value))

	for key, ref := range value {
		newMap[key] = RefToInterface[T, I](ref)
	}

	return newMap
}

func MapOfSlicesOfRefsToInterfaces[K comparable, T any, I any](value map[K][]*T) map[K][]I {
	newMap := make(map[K][]I, len(value))

	for key, refs := range value {
		newMap[key] = SliceOfRefsToInterfaces[T, I](refs)
	}

	return newMap
}

func SliceOfRefsToInterfaces[T any, I any](value []*T) []I {
	result := make([]I, len(value))

	for i, v := range value {
		result[i] = RefToInterface[T, I](v)
	}

	return result
}

func RefToInterface[T any, I any](ref *T) I {
	var vAny any = ref

	return vAny.(I) //nolint:forcetypeassert,errcheck // @todo
}
