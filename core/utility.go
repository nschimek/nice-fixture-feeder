package core

import "time"

var (
	Version = "v0.0.0" // this will get set by -idflags at build time
	EST, _ = time.LoadLocation("America/New_York")
	CST, _ = time.LoadLocation("America/Chicago")
	UTC, _ = time.LoadLocation("UTC")
	Exists = struct{}{}
	YYYY_MM_DD = "2006-01-02"
)

func IdMapToArray(idMap map[int]struct{}) (ids []int) {
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

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