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

// A function that takes in a row and returns the data necessary for selecting it from the table
type RowSelector func([]interface{}) []interface{}

type ICrud interface {
	Create(...interface{}) interface{}
	Retrieve(interface{}) interface{}
	Update(interface{}, interface{}) interface{}
	Delete(interface{})
}

type CrudUpdate struct {
	Key     interface{}
	NewData interface{}
}

type Crud struct {
	TableName         string
	selectorFormatter SelectorFormatter
	toArrFunc         ToArrFunc
	scannerFunc       db.RowScanner
	rowSelector       RowSelector
	fromArrFunc       FromArrFunc
	createFunction    CreateFunc
}
