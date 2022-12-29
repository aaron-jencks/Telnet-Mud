package db

import (
	"encoding/json"
	"fmt"
	"mud/utils/io/csv"
	"mud/utils/ui/logger"
	"os"
	"strings"
)

var DB_LOCATION string = "./data"

func createDbPath(name string, ending string) string {
	return fmt.Sprintf("%s/%s.%s", DB_LOCATION, name, ending)
}

func createCSVPath(name string) string {
	return createDbPath(name, "csv")
}

func createJsonPath(name string) string {
	return createDbPath(name, "json")
}

func TableExists(tableName string) bool {
	_, err := os.Stat(createCSVPath(tableName))
	return !os.IsNotExist(err)
}

func DbDirectoryExists() bool {
	_, err := os.Stat(DB_LOCATION)
	return !os.IsNotExist(err)
}

type TableInfo struct {
	ColumnTypes   []string
	PrimaryKey    int
	UniquePrimary bool
	PrimaryIndex  map[string][]int64
	Indices       map[string]map[string][]int64
}

type TableDefinition struct {
	Name  string
	CSV   *csv.CSVFile
	Info  TableInfo
	Cache *DataCache
}

func checkError(e interface{}) {
	if e != nil {
		logger.ErrorCustomCaller(1, e)
		panic(e)
	}
}

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

func UpdateJson(tableName string, info TableInfo) {
	f, err := os.OpenFile(createJsonPath(tableName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer f.Close()
	checkError(err)
	encoder := json.NewEncoder(f)
	err = encoder.Encode(info)
	checkError(err)
}

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

func DeleteTable(tableName string) {
	if TableExists(tableName) {
		os.Remove(createCSVPath(tableName))
		os.Remove(createJsonPath(tableName))
	}

	logger.Warn("Table %s did not exist", tableName)
}

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

func convertToColumnType(data interface{}, typename string) string {
	conversionString := getConversionString(typename)
	if typename == "string" {
		data = prepareStringForStorage(data.(string))
	}
	result := fmt.Sprintf(conversionString, data)
	return result
}

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

func (td *TableDefinition) DeleteLine(line int64) {
	td.CSV.DeleteLine(line)
	td.Cache.DeleteEntry(line)

	td.UpdateIndices()
}

func (td *TableDefinition) DeleteDataByKey(key interface{}) {
	qdata := td.Query(key, td.CSV.Columns[td.Info.PrimaryKey])
	for _, qline := range qdata {
		td.CSV.DeleteLine(int64(qline[0].(int)))
	}

	td.UpdateIndices()
}

func (td TableDefinition) isIndexed(column string) (bool, bool) {
	_, nIndex := td.Info.Indices[column]

	cindex := stringColumnToInt(column, td.CSV.Columns)

	return nIndex, cindex >= 0 && cindex == td.Info.PrimaryKey
}

func (td *TableDefinition) RetrieveLine(line int64) []interface{} {
	if td.Cache.Exists(line) {
		return td.Cache.RetrieveEntry(line)
	}

	var result []interface{} = make([]interface{}, len(td.CSV.Columns)+1)
	result[0] = int(line)

	values := td.CSV.ReadSpecificLine(line)

	for ci, cvalue := range values {
		result[ci+1] = convertFromString(cvalue, td.Info.ColumnTypes[ci])
	}

	td.Cache.InsertValue(line, result)

	return result
}

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

func (td *TableDefinition) QueryPK(key interface{}) [][]interface{} {
	return td.Query(key, td.CSV.Columns[td.Info.PrimaryKey])
}

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

func (td *TableDefinition) ModifyRowColumn(line int, column string, value interface{}) {
	currentLine := td.RetrieveLine(int64(line))[1:]
	cindex := stringColumnToInt(column, td.CSV.Columns)
	currentLine[cindex] = value

	td.ModifyRow(line, currentLine)
}
