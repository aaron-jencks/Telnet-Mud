package room

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/services/tile"
	"mud/utils/crud"
	"mud/utils/io/db"
)

type ExpandedRoom struct {
	Id             int
	Name           string
	Description    string
	Height         int
	Width          int
	BackgroundTile entities.Tile
}

func roomToArr(rs interface{}) []interface{} {
	re := rs.(entities.Room)
	return []interface{}{
		re.Id,
		re.Name,
		re.Description,
		re.Height,
		re.Width,
		re.BackgroundTile,
	}
}

func roomFromArr(arr []interface{}) interface{} {
	return entities.Room{
		Id:             arr[1].(int),
		Name:           arr[2].(string),
		Description:    arr[3].(string),
		Height:         arr[4].(int),
		Width:          arr[5].(int),
		BackgroundTile: arr[6].(string),
	}
}

func createRoomFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return args
}

func roomSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func roomScanner(row *sql.Rows) (interface{}, error) {
	result := ExpandedRoom{}
	var backgroundTile string
	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Height, &result.Width, &backgroundTile)
	if err != nil {
		return nil, err
	}
	tile := tile.CRUD.Retrieve(backgroundTile)
	result.BackgroundTile = tile.(entities.Tile)
	return result, nil
}

func roomUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(ExpandedRoom)
	nis := newValue.(ExpandedRoom)

	var result []crud.RowModStruct

	if ois.Name != nis.Name {
		result = append(result, crud.RowModStruct{
			Column:   "Name",
			NewValue: fmt.Sprintf("\"%s\"", nis.Name),
		})
	}
	if ois.Description != nis.Description {
		result = append(result, crud.RowModStruct{
			Column:   "Description",
			NewValue: fmt.Sprintf("\"%s\"", nis.Description),
		})
	}
	if ois.Height != nis.Height {
		result = append(result, crud.RowModStruct{
			Column:   "Height",
			NewValue: nis.Height,
		})
	}
	if ois.Width != nis.Width {
		result = append(result, crud.RowModStruct{
			Column:   "Width",
			NewValue: nis.Width,
		})
	}
	if ois.BackgroundTile.Name != nis.BackgroundTile.Name {
		result = append(result, crud.RowModStruct{
			Column:   "BackgroundTile",
			NewValue: fmt.Sprintf("\"%s\"", nis.BackgroundTile.Name),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("rooms", roomSelector, roomToArr, roomScanner, roomFromArr, createRoomFunc, roomUpdateFunc)
