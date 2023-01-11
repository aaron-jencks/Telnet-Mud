package db

import "database/sql"

type RowScanner func(*sql.Rows) (interface{}, error)
