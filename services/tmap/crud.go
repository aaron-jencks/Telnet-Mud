package tmap

import (
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

func createTileFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	result := []interface{}{args[0], args[1], args[2], args[3]}

	if len(args) == 5 {
		result = append(result, args[4])
	} else {
		// Place tile on top
		topZ := -1
		roomTiles := table.QueryPK(args[0])
		for _, tile := range roomTiles {
			if tile[3] == args[2] && tile[4] == args[3] && tile[5].(int) > topZ {
				topZ = tile[5].(int)
			}
		}
		result = append(result, topZ+1)
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("tiles", tileToArr, tileFromArr, createTileFunc)

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
