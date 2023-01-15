package crud

import (
	"database/sql"
	"fmt"
	"mud/utils"
	"mud/utils/io/db"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestTable() {
	db.CreateTableIfNotExist(
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
	db.DeleteTable("TestTable")
	setupTestTable()
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

func TestMain(m *testing.M) {
	// Setup
	resetTable()

	// Run the tests
	ecode := m.Run()

	// Cleanup
	if db.DbDirectoryExists() {
		os.RemoveAll(filepath.Dir(utils.DB_LOCATION))
	}

	os.Exit(ecode)
}

func getTestSelectorFormatter(t *testing.T) SelectorFormatter {
	return func(args []interface{}) string {
		assert.Equal(t, 1, len(args), "selector formatter should be called with a single argument")
		return fmt.Sprintf("RID=%d", args[0].(int64))
	}
}

func testToArrFunc(i interface{}) []interface{} {
	iv := i.(testStruct)
	return []interface{}{
		iv.RId,
		iv.Column1,
		iv.Column2,
		iv.Column3,
	}
}

func getTestFromArrFunc(t *testing.T) FromArrFunc {
	return func(i []interface{}) interface{} {
		assert.Equal(t, 4, len(i), "From Array Func should be called with 4 arguments")
		return testStruct{
			int(i[0].(int64)),
			i[1].(string),
			i[2].(string),
			i[3].(int),
		}
	}
}

func testCreateFunction(table db.TableDefinition, args ...interface{}) []interface{} {
	return args
}

func testUpdateFunc(oldValue, newValue interface{}) []RowModStruct {
	ovs := oldValue.(testStruct)
	nvs := newValue.(testStruct)

	var result []RowModStruct
	if ovs.Column1 != nvs.Column1 {
		result = append(result, RowModStruct{
			Column:   "Column1",
			NewValue: fmt.Sprintf("\"%s\"", nvs.Column1),
		})
	}
	if ovs.Column2 != nvs.Column2 {
		result = append(result, RowModStruct{
			Column:   "Column2",
			NewValue: fmt.Sprintf("\"%s\"", nvs.Column2),
		})
	}
	if ovs.Column3 != nvs.Column3 {
		result = append(result, RowModStruct{
			Column:   "Column3",
			NewValue: nvs.Column3,
		})
	}
	if ovs.RId != nvs.RId {
		result = append(result, RowModStruct{
			Column:   "RId",
			NewValue: nvs.RId,
		})
	}

	return result
}

func getTestCrud(t *testing.T) Crud {
	return CreateCrud(
		"TestTable",
		getTestSelectorFormatter(t),
		testToArrFunc,
		testRowScanner,
		getTestFromArrFunc(t),
		testCreateFunction,
		testUpdateFunc,
	)
}

func TestCreateRetrieve(t *testing.T) {
	resetTable()
	crud := getTestCrud(t)
	nid := crud.Create("TestValue", "Something, else", 1234)
	assert.Equal(t, int64(1), nid, "Should create the correct id")
	nv := crud.Retrieve(nid)
	assert.Equal(t, testStruct{
		1,
		"TestValue",
		"Something, else",
		1234,
	}, nv, "Newly created value should equal the passed in arguments")
}

func TestUpdate(t *testing.T) {
	TestCreateRetrieve(t)
	crud := getTestCrud(t)
	ov := crud.Retrieve(int64(1)).(testStruct)
	ov.Column2 = "Hello Aaron"
	nv := crud.Update(ov, int64(1))
	assert.Equal(t, ov, nv, "Updated value should be returned by update")
}

func TestDelete(t *testing.T) {
	TestCreateRetrieve(t)
	crud := getTestCrud(t)
	crud.Delete(int64(1))
	missingValue := crud.Retrieve(int64(1))
	assert.Nil(t, missingValue, "Should delete data from the table")
}
