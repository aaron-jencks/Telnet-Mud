package item

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func itemToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Item)
	return []interface{}{
		rec.Id,
		rec.Name,
		rec.Description,
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
