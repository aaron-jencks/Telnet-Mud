package detail

import (
	"math/rand"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func detailToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Detail)
	return []interface{}{
		rec.Id,
		rec.Room,
		rec.Direction,
		rec.Detail,
		rec.Perception,
	}
}

func detailFromArr(arr []interface{}) interface{} {
	return entities.Detail{
		Id:         arr[1].(int),
		Room:       arr[2].(int),
		Direction:  arr[3].(string),
		Detail:     arr[4].(string),
		Perception: arr[5].(int),
	}
}

func createDetailFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}
	return []interface{}{id, args[0], args[1], args[2]}
}

var CRUD crud.Crud = crud.CreateCrud("details", detailToArr, detailFromArr, createDetailFunc)

func GetRoomDetails(r entities.Room) []entities.Detail {
	table := CRUD.FetchTable()
	rows := table.Query(r.Id, "Room")

	var result []entities.Detail = make([]entities.Detail, len(rows))
	for ri, row := range rows {
		result[ri] = detailFromArr(row).(entities.Detail)
	}

	return result
}

func TestPerception(min, max int) bool {
	return rand.Intn(max) >= min
}
