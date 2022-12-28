package note

import (
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

func createNoteFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}
	return []interface{}{id, args[0], args[1], args[2]}
}

var CRUD crud.Crud = crud.CreateCrud("notes", noteToArr, noteFromArr, createNoteFunc)
