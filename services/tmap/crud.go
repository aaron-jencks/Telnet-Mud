package tmap

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

var TILENOTFOUND entities.Map = entities.Map{
	Tile: "Not Found",
	Room: -1,
	X:    -1,
	Y:    -1,
	Z:    -1,
}

func tileToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Map)
	return []interface{}{
		rec.Room,
		rec.Tile,
		rec.X,
		rec.Y,
		rec.Z,
	}
}

func tileFromArr(arr []interface{}) interface{} {
	return entities.Map{
		Room: arr[1].(int),
		Tile: arr[2].(string),
		X:    arr[3].(int),
		Y:    arr[4].(int),
		Z:    arr[5].(int),
	}
}

func createTileFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	result := []interface{}{args[0], args[1], args[2], args[3]}

	if len(args) == 5 {
		result = append(result, args[4])
	} else {
		// Place tile on top
		topZ := -1
		roomTiles := table.QueryData(fmt.Sprintf("Room=%d and X=%d and Y=%d", args[0].(int), args[2].(int), args[3].(int)), mapScanner)
		for _, tile := range roomTiles {
			ts := tile.(entities.Map)
			if ts.Z > topZ {
				topZ = ts.Z
			}
		}
		result = append(result, topZ+1)
	}

	return result
}

func mapSelector(args []interface{}) string {
	return fmt.Sprintf("Room=%d and X=%d and Y=%d and Z=%d",
		args[0].(int), args[1].(int), args[2].(int), args[3].(int))
}

func mapScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Map{}
	err := row.Scan(&result.Room, &result.Tile, &result.X, &result.Y, &result.Z)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func mapUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Map)
	nis := newValue.(entities.Map)

	var result []crud.RowModStruct

	if ois.Room != nis.Room {
		result = append(result, crud.RowModStruct{
			Column:   "Room",
			NewValue: nis.Room,
		})
	}
	if ois.Tile != nis.Tile {
		result = append(result, crud.RowModStruct{
			Column:   "Tile",
			NewValue: fmt.Sprintf("\"%s\"", nis.Tile),
		})
	}
	if ois.X != nis.X {
		result = append(result, crud.RowModStruct{
			Column:   "X",
			NewValue: nis.X,
		})
	}
	if ois.Y != nis.Y {
		result = append(result, crud.RowModStruct{
			Column:   "Y",
			NewValue: nis.Y,
		})
	}
	if ois.Z != nis.Z {
		result = append(result, crud.RowModStruct{
			Column:   "Z",
			NewValue: nis.Z,
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("map", mapSelector, tileToArr, mapScanner, tileFromArr, createTileFunc, mapUpdateFunc)

// TODO redo these to make use of queries

func GetTilesForRoom(room int) []entities.Map {
	tiles := CRUD.RetrieveAll(room)

	var result []entities.Map = make([]entities.Map, len(tiles))
	for ti, tile := range tiles {
		result[ti] = tile.(entities.Map)
	}

	return result
}

func GetCurrentTilesForCoord(room, x, y int) []entities.Map {
	roomTiles := GetTilesForRoom(room)

	var result []entities.Map
	for _, tile := range roomTiles {
		if tile.X == x && tile.Y == y {
			result = append(result, tile)
		}
	}

	return result
}

func GetTileForCoord(room, x, y, z int) entities.Map {
	roomTiles := GetTilesForRoom(room)

	for _, tile := range roomTiles {
		if tile.X == x && tile.Y == y && tile.Z == z {
			return tile
		}
	}

	return TILENOTFOUND
}

func GetTopMostTile(room, x, y int) entities.Map {
	tiles := GetCurrentTilesForCoord(room, x, y)

	var maxT entities.Map
	for _, tile := range tiles {
		if tile.Z > maxT.Z {
			maxT = tile
		}
	}

	return maxT
}

func GetTilesForRegion(room, trX, trY, blX, blY int) []entities.Map {
	roomTiles := GetTilesForRoom(room)

	var regionTiles []entities.Map
	for _, tile := range roomTiles {
		if tile.X >= trX && tile.X <= blX && tile.Y >= trY && tile.Y <= blY {
			regionTiles = append(regionTiles, tile)
		}
	}

	return regionTiles
}
