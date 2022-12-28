package inventory

import (
	"mud/entities"
	"mud/services/item"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func inventoryToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Inventory)
	return []interface{}{
		rec.Id,
		rec.Player,
		rec.Item,
		rec.Quantity,
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

type ExpandedInventory struct {
	Item     entities.Item
	Quantity int
}

func GetPlayerInventory(p entities.Player) []ExpandedInventory {
	table := CRUD.FetchTable()
	rows := table.Query(p.Id, "Player")

	var result []ExpandedInventory = make([]ExpandedInventory, len(rows))

	for ri, row := range rows {
		result[ri].Item = item.CRUD.Retrieve(row[3]).(entities.Item)
		result[ri].Quantity = row[4].(int)
	}

	return result
}
