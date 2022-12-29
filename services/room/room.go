package room

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func roomToArr(rs interface{}) []interface{} {
	re := rs.(entities.Room)
	return []interface{}{
		re.Id,
		re.Name,
		re.Description,
		re.Height,
		re.Width,
	}
}

func roomFromArr(arr []interface{}) interface{} {
	return entities.Room{
		Id:          arr[1].(int),
		Name:        arr[2].(string),
		Description: arr[3].(string),
		Height:      arr[4].(int),
		Width:       arr[5].(int),
	}
}

func createRoomFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}

	return []interface{}{id, args[0], args[1], args[2], args[3]}
}

var CRUD crud.Crud = crud.CreateCrud("rooms", roomToArr, roomFromArr, createRoomFunc)
