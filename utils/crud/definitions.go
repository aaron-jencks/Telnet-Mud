package crud

import "mud/utils/io/db"

var tableMap map[string]db.TableDefinition = make(map[string]db.TableDefinition)

// Formats an array of data into a selector string for a WHERE clause
type SelectorFormatter func([]interface{}) string

// Converts an entity into row data
type ToArrFunc func(interface{}) []interface{}

// Converts the given row data into an entity
type FromArrFunc func([]interface{}) interface{}

// A function for creating new rows to insert into the datatable
// Returns the row to insert, and the selector data for that row
type CreateFunc func(db.TableDefinition, ...interface{}) []interface{}

type RowModStruct struct {
	Column   string
	NewValue interface{}
}

type UpdateFunc func(interface{}, interface{}) []RowModStruct

type Crud struct {
	TableName         string
	selectorFormatter SelectorFormatter
	toArrFunc         ToArrFunc
	scannerFunc       db.RowScanner
	fromArrFunc       FromArrFunc
	createFunction    CreateFunc
	updateFunc        UpdateFunc
}
