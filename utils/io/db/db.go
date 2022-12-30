package db

import (
	"encoding/json"
	"fmt"
	"mud/utils/io/csv"
	"mud/utils/ui/logger"
	"os"
	"strings"
)

// The default location of the data tables
// defaults to "./data"
var DB_LOCATION string = "./data"

// Creates the correct path for the given file name and extension type
func createDbPath(name string, ending string) string {
	return fmt.Sprintf("%s/%s.%s", DB_LOCATION, name, ending)
}

// Returns the correct csv file path for the given table name
func createCSVPath(name string) string {
	return createDbPath(name, "csv")
}

// Returns the correct json file path for the given table name
func createJsonPath(name string) string {
	return createDbPath(name, "json")
}

// Determines if the table file exists for the given table name
func TableExists(tableName string) bool {
	_, err := os.Stat(createCSVPath(tableName))
	return !os.IsNotExist(err)
}

// Determines if the DB_LOCATION folder exists
func DbDirectoryExists() bool {
	_, err := os.Stat(DB_LOCATION)
	return !os.IsNotExist(err)
}

// Contains data for the table information
// translation information and indices
type TableInfo struct {
	ColumnTypes   []string                      // The types for the database table, currently ("string", or "integer")
	PrimaryKey    int                           // The column index of the primary key column
	UniquePrimary bool                          // Indicates if the primary column should be unique or not (not currently implemented)
	PrimaryIndex  map[string][]int64            // Indicates which lines correspond to each distinct entry in the primary key column
	Indices       map[string]map[string][]int64 // Indicates which lines corresponse to indices created for other columns
}

// Represents a data table
// Contains the name of the table as well as the csv file
// and a cache  for requests.
type TableDefinition struct {
	Name  string       // the name of the data table
	CSV   *csv.CSVFile // the csv handler for the table
	Info  TableInfo    // contains other important information for the table
	Cache *DataCache   // The in-memory cache for the table
}

func checkError(e interface{}) {
	if e != nil {
		logger.ErrorCustomCaller(1, e)
		panic(e)
	}
}

// Fetches an existing table definition from it's json and csv files
func FetchTableDefinition(tableName string) TableDefinition {
	file := csv.ParseCSV(createCSVPath(tableName))
	var jobj TableInfo

	f, err := os.Open(createJsonPath(tableName))
	defer f.Close()
	checkError(err)

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&jobj)
	checkError(err)

	return TableDefinition{
		tableName,
		&file,
		jobj,
		&DataCache{
			make(map[int64][]interface{}),
			make(map[int64]int64),
		},
	}
}

// Converts a column name to it's corresponding integer index
func stringColumnToInt(column string, columns []string) int {
	var cindex int = -1

	for ci, col := range columns {
		if col == column {
			cindex = int(ci)
			break
		}
	}

	if cindex < 0 {
		logger.Warn("The column %s was not found in the list of %v", column, columns)
	}

	return cindex
}

// Creates a new index for a given column
// An index is a map of unique column entries to their corresponding
// line numbers in the csv file
func CreateIndex(file *csv.CSVFile, column string) map[string][]int64 {
	var index map[string][]int64 = make(map[string][]int64)
	var cindex int = stringColumnToInt(column, file.Columns)
	var li int64

	if cindex < 0 {
		logger.Error("Table %s has no column %s", file.Filepath, column)
		panic("No such Column")
	}

	for li = 0; li < file.LineCount; li++ {
		line := file.ReadSpecificLine(li)

		cv := line[cindex]
		_, present := index[cv]
		if !present {
			index[cv] = []int64{}
		}
		index[cv] = append(index[cv], li)
	}

	return index
}

// Updates the json file with the table information
func UpdateJson(tableName string, info TableInfo) {
	f, err := os.OpenFile(createJsonPath(tableName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	checkError(err)
	encoder := json.NewEncoder(f)
	err = encoder.Encode(info)
	checkError(err)
}

// Either fetches or creates the table
// If fetched, all arguments aside from the tablename are not used.
func CreateTableIfNotExist(tableName string, columns []string,
	columnTypes []string, primaryKey int, uniquePrimary bool) TableDefinition {
	if TableExists(tableName) {
		return FetchTableDefinition(tableName)
	}

	if !DbDirectoryExists() {
		logger.Info("db directory %s did not exist, creating...", DB_LOCATION)
		err := os.MkdirAll(DB_LOCATION, 0777)
		checkError(err)
	}

	file := csv.CreateCSV(createCSVPath(tableName), columns, [][]string{})

	tableInfo := TableInfo{
		columnTypes,
		primaryKey,
		uniquePrimary,
		make(map[string][]int64),
		make(map[string]map[string][]int64),
	}

	UpdateJson(tableName, tableInfo)

	if primaryKey >= 0 {
		tableInfo.PrimaryIndex = CreateIndex(&file, columns[primaryKey])
	} else {
		logger.Warn("No primary key given, no index made")
	}

	return TableDefinition{
		tableName,
		&file,
		tableInfo,
		&DataCache{
			make(map[int64][]interface{}),
			make(map[int64]int64),
		},
	}

}

// Deletes a table if it exists
func DeleteTable(tableName string) {
	if TableExists(tableName) {
		os.Remove(createCSVPath(tableName))
		os.Remove(createJsonPath(tableName))
	}

	logger.Warn("Table %s did not exist", tableName)
}

// Returns a format string for use in fmt.Sprintf
// based on the typename given
// currently supports "string" and "integer"
func getConversionString(typename string) string {
	var conversionString string = "%v"

	switch typename {
	case "string":
		conversionString = "\"%s\""
	case "integer":
		conversionString = "%d"
	}

	return conversionString
}

// Converts a string for storage, this
// involves replacing possibly erroneously parsable data
// with escaped versions
func prepareStringForStorage(data string) string {
	escapeMap := map[string]string{
		" ":  "\\s",
		"\"": "\\\"",
		"\t": "\\t",
		"\n": "\\n",
		"\r": "\\r",
	}

	for k, v := range escapeMap {
		data = strings.ReplaceAll(data, k, v)
	}

	return data
}

// Undoes the conversions performed in prepareStringForStorage
func unpackStoredString(data string) string {
	escapeMap := map[string]string{
		"\\s":  " ",
		"\\\"": "\"",
		"\\t":  "\t",
		"\\n":  "\n",
		"\\r":  "\r",
	}

	for k, v := range escapeMap {
		data = strings.ReplaceAll(data, k, v)
	}

	return data
}

// Converts an any interface type to the corresponding typed string
// for storage in the csv
func convertToColumnType(data interface{}, typename string) string {
	conversionString := getConversionString(typename)
	if typename == "string" {
		data = prepareStringForStorage(data.(string))
	}
	result := fmt.Sprintf(conversionString, data)
	return result
}

// Converts a stored string in the csv back into the correct typed data.
func convertFromString(sdata, typename string) interface{} {
	conversionString := getConversionString(typename)

	switch typename {
	case "string":
		var result string
		_, err := fmt.Sscanf(sdata[1:len(sdata)-1], "%s", &result)
		checkError(err)
		return unpackStoredString(result)
	case "integer":
		var result int
		_, err := fmt.Sscanf(sdata, conversionString, &result)
		checkError(err)
		return result
	default:
		logger.Error("Unknown conversion type %s for value %s", typename, sdata)
		panic("Unknown conversion type")
	}
}

// Reparses indices, this is performance heavy,
// but necessary when modifying the database
func (td *TableDefinition) UpdateIndices() {
	// TODO usages of this could be more efficient

	var indexedColumns []string

	for col := range td.Info.Indices {
		indexedColumns = append(indexedColumns, col)
	}

	td.Info.Indices = make(map[string]map[string][]int64)

	for _, col := range indexedColumns {
		td.Info.Indices[col] = CreateIndex(td.CSV, col)
	}

	if td.Info.PrimaryKey >= 0 {
		td.Info.PrimaryIndex = CreateIndex(td.CSV, td.CSV.Columns[td.Info.PrimaryKey])
	}

	UpdateJson(td.Name, td.Info)
}

// Adds new data to the data table and returns the number of lines added
func (td *TableDefinition) AddData(data [][]interface{}) int {
	for _, dline := range data {
		var sdata []string = make([]string, len(dline))
		for sdi, col := range dline {
			sdata[sdi] = convertToColumnType(col, td.Info.ColumnTypes[sdi])
		}
		td.CSV.AppendLine(sdata)

		// Updates the indices
		line := td.CSV.LineCount - 1
		if td.Info.PrimaryKey >= 0 {
			pv := sdata[td.Info.PrimaryKey]
			_, exists := td.Info.PrimaryIndex[pv]
			if !exists {
				td.Info.PrimaryIndex[pv] = []int64{}
			}
			td.Info.PrimaryIndex[pv] = append(td.Info.PrimaryIndex[pv], line)
		}
		td.Cache.InsertValue(line, dline)

		for k := range td.Info.Indices {
			cindex := stringColumnToInt(k, td.CSV.Columns)
			v := sdata[cindex]
			arr, exists := td.Info.Indices[k][v]
			if !exists {
				arr = []int64{}
			}
			td.Info.Indices[k][v] = append(arr, line)
		}
	}

	UpdateJson(td.Name, td.Info)

	return len(data)
}

// Deletes a single line from the data table
func (td *TableDefinition) DeleteLine(line int64) {
	defer td.UpdateIndices()

	td.CSV.DeleteLine(line)
	td.Cache.DeleteEntry(line)
}

// Deletes multiple lines from the data table
func (td *TableDefinition) DeleteLines(lines []int64) {
	defer td.UpdateIndices()

	td.CSV.DeleteLines(lines)

	for _, line := range lines {
		td.Cache.DeleteEntry(line)
	}
}

// Deletes all lines matching the given primary key
func (td *TableDefinition) DeleteDataByKey(key interface{}) {
	qdata := td.Query(key, td.CSV.Columns[td.Info.PrimaryKey])

	var lines []int64
	for _, qline := range qdata {
		lines = append(lines, int64(qline[0].(int)))
	}

	td.DeleteLines(lines)
}

// Returns whether the given column is indexed or not
func (td TableDefinition) isIndexed(column string) (bool, bool) {
	_, nIndex := td.Info.Indices[column]

	cindex := stringColumnToInt(column, td.CSV.Columns)

	return nIndex, cindex >= 0 && cindex == td.Info.PrimaryKey
}

// Retrieves a specific line from the data table
func (td *TableDefinition) RetrieveLine(line int64) []interface{} {
	if td.Cache.Exists(line) {
		lineData := td.Cache.RetrieveEntry(line)
		return append([]interface{}{int(line)}, lineData...)
	}

	var result []interface{} = make([]interface{}, len(td.CSV.Columns)+1)
	result[0] = int(line)

	values := td.CSV.ReadSpecificLine(line)

	for ci, cvalue := range values {
		result[ci+1] = convertFromString(cvalue, td.Info.ColumnTypes[ci])
	}

	td.Cache.InsertValue(line, result[1:])

	return result
}

// Query's the table for the given value in the given column
// Returns all rows matching that value
func (td *TableDefinition) Query(value interface{}, column string) [][]interface{} {
	var results [][]interface{}

	cindex := stringColumnToInt(column, td.CSV.Columns)
	normIndex, primIndex := td.isIndexed(column)
	svalue := convertToColumnType(value, td.Info.ColumnTypes[cindex])

	if primIndex {
		// Is the primary key
		v, exists := td.Info.PrimaryIndex[svalue]
		if exists {
			for _, lindex := range v {
				results = append(results, td.RetrieveLine(lindex))
			}
		}

	} else if normIndex {
		// Is an indexed column
		v, exists := td.Info.Indices[column][svalue]
		if exists {
			for _, lindex := range v {
				results = append(results, td.RetrieveLine(lindex))
			}
		}

	} else {
		// Need to traverse the entire file
		var li int64
		for li = 0; li < td.CSV.LineCount; li++ {
			resultLine := td.RetrieveLine(li)
			if resultLine[cindex+1] == value {
				results = append(results, resultLine)
			}
		}
	}

	return results
}

// Queries multiple columns similar to Query
// must pass arguments in in pairs of (value, column)
// where value is the value to look for and column is the column string
func (td *TableDefinition) MultiQuery(args ...interface{}) [][]interface{} {
	var results [][]interface{}

	var pbuffer [][]interface{}
	var buffer [][]interface{}
	for ai := 0; ai < len(args); ai += 2 {
		buffer = td.Query(args[ai], args[ai+1].(string))
		if pbuffer == nil {
			pbuffer = buffer
			continue
		} else {
			// refine results
			results = nil
			for _, brow := range buffer {
				for _, prow := range pbuffer {
					if brow[0] == prow[0] {
						results = append(results, brow)
					}
				}
			}
			pbuffer = results
		}

		if len(pbuffer) == 0 {
			// There are no results that satisfy all of the columns
			return pbuffer
		}
	}

	return results
}

// Queries the primary key column for the given value
func (td *TableDefinition) QueryPK(key interface{}) [][]interface{} {
	return td.Query(key, td.CSV.Columns[td.Info.PrimaryKey])
}

// Updates a specific line in the data table to the new value
func (td *TableDefinition) ModifyRow(line int, data []interface{}) {
	var csvLine []string = make([]string, len(data))

	for di, dv := range data {
		csvLine[di] = convertToColumnType(dv, td.Info.ColumnTypes[di])
	}

	td.CSV.ModifyLine(line, csvLine)

	if !td.Cache.Exists(int64(line)) {
		td.Cache.InsertValue(int64(line), data)
	} else {
		td.Cache.UpdateEntry(int64(line), data)
	}

	td.UpdateIndices()
}

// Modifies a specific column of the given line
func (td *TableDefinition) ModifyRowColumn(line int, column string, value interface{}) {
	currentLine := td.RetrieveLine(int64(line))[1:]
	cindex := stringColumnToInt(column, td.CSV.Columns)
	currentLine[cindex] = value

	td.ModifyRow(line, currentLine)
}
