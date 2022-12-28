package loot

import (
	"mud/entities"
	"mud/services/item"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func lootToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Loot)
	return []interface{}{
		rec.Id,
		rec.Room,
		rec.Item,
		rec.Quantity,
	}
}

func lootFromArr(arr []interface{}) interface{} {
	return entities.Loot{
		Id:       arr[1].(int),
		Room:     arr[2].(int),
		Item:     arr[3].(int),
		Quantity: arr[4].(int),
	}
}

func createLootFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}
	return []interface{}{id, args[0], args[1], args[2]}
}

var CRUD crud.Crud = crud.CreateCrud("loot", lootToArr, lootFromArr, createLootFunc)

type ExpandedLoot struct {
	Id       int
	Room     entities.Room
	Item     entities.Item
	Quantity int
}

func GetLootForRoom(r entities.Room) []ExpandedLoot {
	table := CRUD.FetchTable()
	rows := table.Query(r.Id, "Room")

	var loots []ExpandedLoot = make([]ExpandedLoot, len(rows))
	for ri, row := range rows {
		loots[ri].Id = row[1].(int)
		loots[ri].Room = r

		item := item.CRUD.Retrieve(row[3]).(entities.Item)
		loots[ri].Item = item
		loots[ri].Quantity = row[4].(int)
	}

	return loots
}
