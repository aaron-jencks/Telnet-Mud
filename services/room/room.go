package room_service

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func roomToArr(rs map[string]interface{}) []interface{} {
	return []interface{}{
		int(rs["Id"].(float64)),
		rs["Name"],
		rs["Description"],
	}
}

func roomFromArr(arr []interface{}) interface{} {
	return entities.Room{
		Id:          arr[1].(int),
		Name:        arr[2].(string),
		Description: arr[3].(string),
	}
}

func createRoomFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}

	return []interface{}{id, args[0], args[1]}
}

var CRUD crud.Crud = crud.CreateCrud("rooms", roomToArr, roomFromArr, createRoomFunc)