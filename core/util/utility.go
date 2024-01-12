package util

var Exists = struct{}{}

// IdMapToArray converts a map of ID integers into an array.
func IdMapToArray(idMap map[int]struct{}) (ids []int) {
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

// IdArrayToMap converts an array of integers into a map of integers as they key, with an empty struct value.
func IdArrayToMap(ids []int) (idMap map[int]struct{}) {
	idMap = make(map[int]struct{})

	for _, id := range ids {
		idMap[id] = Exists
	}

	return
}

func MapToArray[K comparable, V any](m map[K]V) (values []V) {
	for _, v := range m {
		values = append(values, v)
	}
	return
}

func MapToKeyArray[K comparable, V any](m map[K]V) (keys []K) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}
