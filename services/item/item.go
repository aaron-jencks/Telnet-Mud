package item

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func itemToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Item)
	return []interface{}{
		rec.Id,
		rec.Name,
		rec.Description,
	}
}

func itemFromArr(arr []interface{}) interface{} {
	return entities.Item{
		Id:          arr[1].(int),
		Name:        arr[2].(string),
		Description: arr[3].(string),
	}
}

func createItemFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1]}
}

func itemSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func itemScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Item{}
	err := row.Scan(&result.Id, &result.Name, &result.Description)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func itemUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Item)
	nis := newValue.(entities.Item)

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

	return result
}

var CRUD crud.Crud = crud.CreateCrud("items", itemSelector, itemToArr, itemScanner, itemFromArr, createItemFunc, itemUpdateFunc)
