package note

import (
	"database/sql"
	"fmt"
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func noteToArr(rs interface{}) []interface{} {
	rec := rs.(entities.Note)
	return []interface{}{
		rec.Id,
		rec.Player,
		rec.Title,
		rec.Contents,
	}
}

func noteFromArr(arr []interface{}) interface{} {
	return entities.Note{
		Id:       arr[1].(int),
		Player:   arr[2].(int),
		Title:    arr[3].(string),
		Contents: arr[4].(string),
	}
}

func createNoteFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1], args[2]}
}

func noteSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d", args[0].(int))
}

func noteScanner(row *sql.Rows) (interface{}, error) {
	result := entities.Note{}
	err := row.Scan(&result.Id, &result.Player, &result.Title, &result.Contents)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func noteUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.Note)
	nis := newValue.(entities.Note)

	var result []crud.RowModStruct

	if ois.Player != nis.Player {
		result = append(result, crud.RowModStruct{
			Column:   "Player",
			NewValue: nis.Player,
		})
	}
	if ois.Title != nis.Title {
		result = append(result, crud.RowModStruct{
			Column:   "Title",
			NewValue: fmt.Sprintf("\"%s\"", nis.Title),
		})
	}
	if ois.Contents != nis.Contents {
		result = append(result, crud.RowModStruct{
			Column:   "Contents",
			NewValue: fmt.Sprintf("\"%s\"", nis.Contents),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("notes", noteSelector, noteToArr, noteScanner, noteFromArr, createNoteFunc, noteUpdateFunc)
