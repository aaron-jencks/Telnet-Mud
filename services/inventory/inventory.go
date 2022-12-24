package inventory_service

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func inventoryToArr(rs map[string]interface{}) []interface{} {
	return []interface{}{
		int(rs["Id"].(float64)),
		int(rs["Player"].(float64)),
		int(rs["Item"].(float64)),
		int(rs["Quantity"].(float64)),
	}
}

func inventoryFromArr(arr []interface{}) interface{} {
	return entities.Inventory{
		Id:       arr[1].(int),
		Player:   arr[2].(int),
		Item:     arr[3].(int),
		Quantity: arr[4].(int),
	}
}

func createInventoryFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}

	return []interface{}{id, args[0], args[1], args[2]}
}

var CRUD crud.Crud = crud.CreateCrud("inventory", inventoryToArr, inventoryFromArr, createInventoryFunc)
