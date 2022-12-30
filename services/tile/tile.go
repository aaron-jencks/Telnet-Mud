package tile

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func tileToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Tile)
	return []interface{}{
		rec.Name,
		rec.IconType,
		rec.Icon,
		rec.BG,
		rec.FG,
	}
}

func tileFromArr(arr []interface{}) interface{} {
	return entities.Tile{
		Name:     arr[1].(string),
		IconType: arr[2].(string),
		Icon:     arr[3].(string),
		BG:       arr[4].(int),
		FG:       arr[5].(int),
	}
}

func createTileFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	result := []interface{}{args[0], args[1], args[2]}

	if len(result) == 5 {
		result = append(result, args[3], args[4])
	} else {
		result = append(result, 0, 30)
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("tiles", tileToArr, tileFromArr, createTileFunc)
