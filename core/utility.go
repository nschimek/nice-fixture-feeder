package core

import "time"

var (
	EST, _ = time.LoadLocation("merica/New_York")
	CST, _ = time.LoadLocation("America/Chicago")
	Exists = struct{}{}
	YYYY_MM_DD = "2006-01-02"
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
		idMap[id] = Exists
	}

	return
}