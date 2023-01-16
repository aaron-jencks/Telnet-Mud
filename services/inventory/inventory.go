package inventory

import (
	"database/sql"
	"fmt"
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

func createInventoryFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1], args[2]}
}

func inventoryScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Inventory{}
	err := row.Scan(&result.Id, &result.Player, &result.Item, &result.Quantity)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func inventorySelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func inventoryUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Inventory)
	nis := newValue.(entities.Inventory)

	var result []crud.RowModStruct

	if ois.Item != nis.Item {
		result = append(result, crud.RowModStruct{
			Column:   "Item",
			NewValue: nis.Item,
		})
	}
	if ois.Player != nis.Player {
		result = append(result, crud.RowModStruct{
			Column:   "Player",
			NewValue: nis.Player,
		})
	}
	if ois.Quantity != nis.Quantity {
		result = append(result, crud.RowModStruct{
			Column:   "Quantity",
			NewValue: nis.Quantity,
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("inventory", inventorySelector,
	inventoryToArr, inventoryScanner, inventoryFromArr,
	createInventoryFunc, inventoryUpdateFunc)

type ExpandedInventory struct {
	Item     entities.Item
	Quantity int
}

func GetPlayerInventory(p entities.Player) []ExpandedInventory {
	table := CRUD.FetchTable()
	rows := table.QueryData(fmt.Sprintf("Player=%d", p.Id), inventoryScanner)

	var result []ExpandedInventory = make([]ExpandedInventory, len(rows))

	for ri, row := range rows {
		rs := row.(entities.Inventory)
		result[ri].Item = item.CRUD.Retrieve(rs.Item).(entities.Item)
		result[ri].Quantity = rs.Quantity
	}

	return result
}

func AddItemToInventory(p entities.Player, i entities.Item, qty int) int {
	table := CRUD.FetchTable()
	rows := table.QueryData(fmt.Sprintf("Player=%d", p.Id), inventoryScanner)

	for _, row := range rows {
		rs := row.(entities.Inventory)
		if rs.Item == i.Id {
			// We already have some
			invent := CRUD.Retrieve(rs.Id).(entities.Inventory)
			invent.Quantity += qty
			CRUD.Update(invent, invent.Id)

			return invent.Quantity
		}
	}

	CRUD.Create(p.Id, i.Id, qty)

	return qty
}
