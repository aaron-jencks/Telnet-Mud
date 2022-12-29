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
	}
}

func tileFromArr(arr []interface{}) interface{} {
	return entities.Tile{
		Name:     arr[1].(string),
		IconType: arr[2].(string),
		Icon:     arr[3].(string),
	}
}

func createTileFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1], args[2]}
}

var CRUD crud.Crud = crud.CreateCrud("tiles", tileToArr, tileFromArr, createTileFunc)
