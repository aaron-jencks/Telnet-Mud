package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"mud/utils"
	mtesting "mud/utils/testing"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestTable() {
	CreateTableIfNotExist(
		"TestTable",
		[]string{
			"RId",
			"Column1",
			"Column2",
			"Column3",
		},
		[]string{
			"RId integer primary key autoincrement",
			"Column1 text",
			"Column2 text",
			"Column3 integer",
		},
		true,
	)
}

func resetTable() {
	DeleteTable("TestTable")
	setupTestTable()
}

func TestMain(m *testing.M) {
	// Setup
	resetTable()

	// Run the tests
	ecode := m.Run()

	// Cleanup
	if DbDirectoryExists() {
		os.RemoveAll(filepath.Dir(utils.DB_LOCATION))
	}

	os.Exit(ecode)
}

func TestTableCreateDelete(t *testing.T) {
	setupTestTable()
	rand.Seed(time.Now().Unix())

	DeleteTable("TestTable")

	setupTestTable()
}

func generateRandomInsertData(n int) [][]interface{} {
	var result [][]interface{} = make([][]interface{}, n)

	for ri := range result {
		result[ri] = []interface{}{
			mtesting.GenerateRandomAlnumString(rand.Intn(10)),
			mtesting.GenerateRandomAlnumString(rand.Intn(10)),
			rand.Int(),
		}
	}

	return result
}

func TestAddDeleteData(t *testing.T) {
	resetTable()
	tbl := FetchTableDefinition("TestTable")
	ids, nmod := tbl.AddData(generateRandomInsertData(10))
	assert.Equal(t, nmod, int64(10), "Adding data should return the number of modified rows")
	for iid, id := range ids {
		assert.Equal(t, int64(iid+1), id, "Autoincrement should work")
		drow := tbl.DeleteData(fmt.Sprintf("RId=%d", id))
		assert.Equal(t, int64(1), drow, "Delete should delete the correct number of rows")
	}

}

type testStruct struct {
	RId     int
	Column1 string
	Column2 string
	Column3 int
}

func testRowScanner(row *sql.Rows) (interface{}, error) {
	var (
		id int
		c1 string
		c2 string
		c3 int
	)
	err := row.Scan(&id, &c1, &c2, &c3)
	if err != nil {
		return nil, err
	}
	return testStruct{
		RId:     id,
		Column1: c1,
		Column2: c2,
		Column3: c3,
	}, nil
}

func TestQueryData(t *testing.T) {
	resetTable()
	tbl := FetchTableDefinition("TestTable")
	rowData := generateRandomInsertData(10)
	ids, _ := tbl.AddData(rowData)
	for _, id := range ids {
		rows := tbl.QueryData(fmt.Sprintf("RId=%d", id), testRowScanner)
		assert.Equal(t, 1, len(rows), "Querying by the pk should return a single row")
		row := rows[0].(testStruct)
		assert.Equal(t, testStruct{
			int(id),
			rowData[id-1][0].(string),
			rowData[id-1][1].(string),
			rowData[id-1][2].(int),
		}, row, "Queried data should equal the inserted data")
		tbl.DeleteData(fmt.Sprintf("RId=%d", id))
	}
}

func TestUpdate(t *testing.T) {
	resetTable()
	tbl := FetchTableDefinition("TestTable")
	rowData := generateRandomInsertData(1)
	ids, _ := tbl.AddData(rowData)
	newColumnValue := rowData[0][2].(int) + 1
	rowCount := tbl.UpdateData(fmt.Sprintf("RId=%d", ids[0]), "Column3", newColumnValue)
	assert.Equal(t, int64(1), rowCount, "Should update the correct number of rows")
	rows := tbl.QueryData(fmt.Sprintf("RId=%d", ids[0]), testRowScanner)
	assert.Equal(t, newColumnValue, rows[0].(testStruct).Column3, "Column should be updated")
}
