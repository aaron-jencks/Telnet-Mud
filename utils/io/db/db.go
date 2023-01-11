package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mud/utils"
	"mud/utils/ui/logger"
	"os"
	"path"
	"strings"
)

func openDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", utils.DB_LOCATION)
	if err != nil {
		panic(err)
	}
	return db
}

func RunQuery(statement string, args ...interface{}) (*sql.Rows, error) {
	db := openDbConnection()
	defer db.Close()
	return db.Query(statement, args...)
}

func RunExec(statement string, args ...interface{}) (sql.Result, error) {
	db := openDbConnection()
	defer db.Close()
	return db.Exec(statement, args...)
}

func RunInsert(statement string, rows [][]interface{}, args ...interface{}) ([]sql.Result, error) {
	db := openDbConnection()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var results []sql.Result
	for _, row := range rows {
		rresult, err := stmt.Exec(row...)
		if err != nil {
			return results, err
		}
		results = append(results, rresult)
	}

	err = tx.Commit()
	if err != nil {
		return results, err
	}

	return results, nil
}

// Determines if the DB_LOCATION folder exists
func DbDirectoryExists() bool {
	_, err := os.Stat(path.Dir(utils.DB_LOCATION))
	return !os.IsNotExist(err)
}

// Represents a data table
// Contains the name of the table as well as the csv file
// and a cache  for requests.
type TableDefinition struct {
	Name        string   // the name of the data table
	ColumnNames []string // The names of the columns in the table
}

func checkError(e interface{}) {
	if e != nil {
		logger.ErrorCustomCaller(1, e)
		panic(e)
	}
}

// Either fetches or creates the table
// If fetched, all arguments aside from the tablename are not used.
func CreateTableIfNotExist(tableName string, columns, columnSpecs []string) TableDefinition {
	if !DbDirectoryExists() {
		logger.Info("db directory %s did not exist, creating...", utils.DB_LOCATION)
		err := os.MkdirAll(utils.DB_LOCATION, 0777)
		checkError(err)
	}

	statement := fmt.Sprintf(`create table if not exists %s (%s)`,
		tableName, strings.Join(columnSpecs, ","))

	_, err := RunExec(statement)
	checkError(err)

	table := TableDefinition{
		tableName,
		columns,
	}

	table.UpdateJson()

	return table
}

// Updates the json file with the table information
func (td TableDefinition) UpdateJson() {
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.json", path.Dir(utils.DB_LOCATION), td.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	checkError(err)
	encoder := json.NewEncoder(f)
	err = encoder.Encode(td)
	checkError(err)
}

// Deletes a table if it exists
func DeleteTable(tableName string) {
	statement := fmt.Sprintf(`delete table if exists %s`, tableName)

	_, err := RunExec(statement)
	checkError(err)

	logger.Warn("Table %s deleted", tableName)
}

// Adds new data to the data table and returns the number of lines added
func (td TableDefinition) AddData(data [][]interface{}) int {
	_, err := RunInsert(fmt.Sprintf(`insert into %s (%s) values (%s)`,
		td.Name, strings.Join(td.ColumnNames, ","),
		strings.Join(strings.Split(strings.Repeat("?", len(td.ColumnNames)), ""), ",")), data)
	checkError(err)

	return len(data)
}

func (td TableDefinition) DeleteData(selector string) int64 {
	result, err := RunExec(fmt.Sprintf("delete from %s where %s", td.Name, selector))
	checkError(err)

	rows, err := result.RowsAffected()
	checkError(err)

	return rows
}

func (td TableDefinition) QueryData(statement string, scanner RowScanner) []interface{} {
	rows, err := RunQuery(statement)
	checkError(err)
	defer rows.Close()

	var result []interface{}
	for rows.Next() {
		value, err := scanner(rows)
		checkError(err)
		result = append(result, value)
	}

	err = rows.Err()
	checkError(err)

	return result
}

func (td TableDefinition) UpdateData(selector, columns string, replacement interface{}) int64 {
	result, err := RunExec(fmt.Sprintf("update %s set %s=%v where %s",
		td.Name, columns, replacement, selector))
	checkError(err)

	rows, err := result.RowsAffected()
	checkError(err)

	return rows
}
