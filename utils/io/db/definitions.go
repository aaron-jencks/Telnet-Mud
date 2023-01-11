package db

import "database/sql"

type RowScanner func(*sql.Rows) (interface{}, error)

// Represents a data table
// Contains the name of the table as well as the csv file
// and a cache  for requests.
type TableDefinition struct {
	Name        string   // the name of the data table
	ColumnNames []string // The names of the columns in the table
	ColumnSpecs []string // The column definitions including the types and any foreign keys/constraints
}
