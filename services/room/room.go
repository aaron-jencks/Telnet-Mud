package room

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func roomToArr(rs interface{}) []interface{} {
	re := rs.(entities.Room)
	return []interface{}{
		re.Id,
		re.Name,
		re.Description,
		re.Height,
		re.Width,
	}
}

func roomFromArr(arr []interface{}) interface{} {
	return entities.Room{
		Id:          arr[1].(int),
		Name:        arr[2].(string),
		Description: arr[3].(string),
		Height:      arr[4].(int),
		Width:       arr[5].(int),
	}
}

func createRoomFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1], args[2], args[3]}
}

func roomSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func roomScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Room{}
	err := row.Scan(&result.Id, &result.Name, &result.Description, &result.Height, &result.Width)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func roomUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Room)
	nis := newValue.(entities.Room)

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

	return result
}

var CRUD crud.Crud = crud.CreateCrud("rooms", roomSelector, roomToArr, roomScanner, roomFromArr, createRoomFunc, roomUpdateFunc)
