package item_service

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func itemToArr(rs map[string]interface{}) []interface{} {
	return []interface{}{
		int(rs["Id"].(float64)),
		rs["Name"],
		rs["Description"],
	}
}

func itemFromArr(arr []interface{}) interface{} {
	return entities.Item{
		Id:          arr[1].(int),
		Name:        arr[2].(string),
		Description: arr[3].(string),
	}
}

func createItemFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}
	return []interface{}{id, args[0], args[1]}
}

var CRUD crud.Crud = crud.CreateCrud("items", itemToArr, itemFromArr, createItemFunc)
