package loot

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/services/item"
	"mud/services/room"
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

func createLootFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	result := []interface{}{}
	result = append(result, args...)

	return result
}

func lootSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func lootScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Loot{}
	err := row.Scan(&result.Id, &result.Room, &result.Item, &result.Quantity,
		&result.X, &result.Y, &result.Z)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func lootUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ols := oldValue.(entities.Loot)
	nls := newValue.(entities.Loot)

	var result []crud.RowModStruct

	if ols.Room != nls.Room {
		result = append(result, crud.RowModStruct{
			Column:   "Room",
			NewValue: nls.Room,
		})
	}
	if ols.Item != nls.Item {
		result = append(result, crud.RowModStruct{
			Column:   "Item",
			NewValue: nls.Item,
		})
	}
	if ols.Quantity != nls.Quantity {
		result = append(result, crud.RowModStruct{
			Column:   "Quantity",
			NewValue: nls.Quantity,
		})
	}
	if ols.X != nls.X {
		result = append(result, crud.RowModStruct{
			Column:   "X",
			NewValue: nls.X,
		})
	}
	if ols.Y != nls.Y {
		result = append(result, crud.RowModStruct{
			Column:   "Y",
			NewValue: nls.Y,
		})
	}
	if ols.Z != nls.Z {
		result = append(result, crud.RowModStruct{
			Column:   "Z",
			NewValue: nls.Z,
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("loot", lootSelector, lootToArr, lootScanner, lootFromArr, createLootFunc, lootUpdateFunc)

type ExpandedLoot struct {
	Id       int
	Room     room.ExpandedRoom
	Item     entities.Item
	Quantity int
	X        int
	Y        int
	Z        int
}

func GetLootForRoom(r room.ExpandedRoom) []ExpandedLoot {
	table := CRUD.FetchTable()
	rows := table.QueryData(fmt.Sprintf("Room=%d", r.Id), lootScanner)

	var loots []ExpandedLoot = make([]ExpandedLoot, len(rows))
	for ri, row := range rows {
		rs := row.(entities.Loot)
		loots[ri].Id = rs.Id
		loots[ri].Room = r

		item := item.CRUD.Retrieve(rs.Item).(entities.Item)
		loots[ri].Item = item
		loots[ri].Quantity = rs.Quantity
		loots[ri].X = rs.X
		loots[ri].Y = rs.Y
		loots[ri].Z = rs.Z
	}

	return loots
}

func GetLootForPosition(r room.ExpandedRoom, x, y int) []ExpandedLoot {
	table := CRUD.FetchTable()
	rows := table.QueryData(fmt.Sprintf("Room=%d and X=%d and Y=%d order by Z desc", r.Id, x, y), lootScanner)

	var loots []ExpandedLoot = make([]ExpandedLoot, len(rows))
	for ri, row := range rows {
		rs := row.(entities.Loot)
		loots[ri].Id = rs.Id
		loots[ri].Room = r

		item := item.CRUD.Retrieve(rs.Item).(entities.Item)
		loots[ri].Item = item
		loots[ri].Quantity = rs.Quantity
		loots[ri].X = rs.X
		loots[ri].Y = rs.Y
		loots[ri].Z = rs.Z
	}

	return loots
}
