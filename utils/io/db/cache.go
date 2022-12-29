package db

import (
	"mud/utils"
)

type DataCache struct {
	Data map[int64][]interface{}
	Uses map[int64]int64
}

func (dc *DataCache) PruneData() {
	var leastUsed int64
	var leastKey int64

	first := true
	for key := range dc.Uses {
		if first {
			leastUsed = dc.Uses[key]
			leastKey = key
			first = false
		} else if dc.Uses[key] < leastUsed {
			leastUsed = dc.Uses[key]
			leastKey = key
		}
	}

	delete(dc.Data, leastKey)
	delete(dc.Uses, leastKey)
}

func (dc *DataCache) InsertValue(row int64, data []interface{}) {
	if dc.Exists(row) {
		return
	}

	if len(dc.Data) == utils.CACHE_SIZE_LIMIT {
		dc.PruneData()
	}

	dc.Data[row] = data
	dc.Uses[row] = 0
}

func (dc *DataCache) RetrieveEntry(row int64) []interface{} {
	v, exists := dc.Data[row]
	if exists {
		dc.Uses[row]++
		return v
	}
	return nil
}

func (dc *DataCache) Exists(row int64) bool {
	_, exists := dc.Data[row]
	return exists
}

func (dc *DataCache) UpdateEntry(row int64, newValue []interface{}) {
	if dc.Exists(row) {
		dc.Data[row] = newValue
		dc.Uses[row]++
	}
}

func (dc *DataCache) DeleteEntry(row int64) {
	if dc.Exists(row) {
		delete(dc.Data, row)
		delete(dc.Uses, row)
	}
}
