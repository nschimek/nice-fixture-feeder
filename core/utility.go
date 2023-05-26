package core

import "time"

var (
	EST, _ = time.LoadLocation("merica/New_York")
	CST, _ = time.LoadLocation("America/Chicago")
	exists = struct{}{}
)

func IdMapToArray(idMap map[string]struct{}) (ids []string) {
	for id := range idMap {
		ids = append(ids, id)
	}
	return
}

func IdArrayToMap(ids []string) (idMap map[string]struct{}) {
	idMap = make(map[string]struct{})

	for _, id := range ids {
		idMap[id] = exists
	}

	return
}