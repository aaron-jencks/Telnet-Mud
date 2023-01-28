package tile

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
	"strings"
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

func createTileFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	result := []interface{}{args[0], args[1], args[2]}

	if len(args) >= 5 {
		result = append(result, args[3:]...)
		if len(args) == 5 {
			result = append(result, false)
		}
	} else {
		result = append(result, 0, 30)
		if len(args) == 4 {
			result = append(result, args[3])
		} else {
			result = append(result, false)
		}
	}

	return result
}

func tileSelector(args []interface{}) string {
	return fmt.Sprintf("Name=\"%s\"", args[0].(string))
}

func tileScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Tile{}
	err := row.Scan(&result.Name, &result.IconType, &result.Icon, &result.BG, &result.FG, &result.Solid)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func tileUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Tile)
	nis := newValue.(entities.Tile)

	var result []crud.RowModStruct

	if ois.Name != nis.Name {
		result = append(result, crud.RowModStruct{
			Column:   "Name",
			NewValue: fmt.Sprintf("\"%s\"", nis.Name),
		})
	}
	if ois.IconType != nis.IconType {
		result = append(result, crud.RowModStruct{
			Column:   "IconType",
			NewValue: fmt.Sprintf("\"%s\"", nis.IconType),
		})
	}
	if ois.Icon != nis.Icon {
		result = append(result, crud.RowModStruct{
			Column:   "Icon",
			NewValue: fmt.Sprintf("\"%s\"", nis.Icon),
		})
	}
	if ois.BG != nis.BG {
		result = append(result, crud.RowModStruct{
			Column:   "BG",
			NewValue: nis.BG,
		})
	}
	if ois.FG != nis.FG {
		result = append(result, crud.RowModStruct{
			Column:   "FG",
			NewValue: nis.FG,
		})
	}
	if ois.Solid != nis.Solid {
		result = append(result, crud.RowModStruct{
			Column:   "Solid",
			NewValue: strings.ToUpper(fmt.Sprintf("%t", nis.Solid)),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("tiles", tileSelector, tileToArr, tileScanner, tileFromArr, createTileFunc, tileUpdateFunc)
