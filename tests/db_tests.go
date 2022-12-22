package db_tests

import (
  "logger"
  "test_defs"
  "assert"
  "db"
  "os"
)

var TEST_DB_TABLE_NAME string = "TestTable"
var TEST_DB_TABLE_COLUMNS []string = []string{
  "Test1",
  "Test2",
  "Test3",
}
var TEST_DB_TABLE_TYPES []string = []string{
  "string",
  "integer",
  "string",
}
var TEST_DB_PK int = 1
var TEST_DB_UNIQUE_PK bool = false
var TEST_DB_LOCATION string = "/tmp/testDbLocation"

var TEST_SUITE test_defs.TestSuite = test_defs.TestSuite {
  "DB API",
  []func() {
    beforeAll,
    TestCreateDelete,
    TestFileParsing,
    TestInsertion,
    TestDeletion,
    TestModification,
    afterAll,
  },
}

func beforeAll() {
  err := os.Mkdir(TEST_DB_LOCATION, 0755)
  if err != nil {
    logger.Error("Error while creating %s: %v", TEST_DB_LOCATION, err)
    panic(err)
  }
}

func createTestData() db.TableDefinition {
  return db.CreateTableIfNotExist(
    TEST_DB_TABLE_NAME,
    TEST_DB_TABLE_COLUMNS,
    TEST_DB_TABLE_TYPES,
    TEST_DB_PK,
    TEST_DB_UNIQUE_PK)
}

func fetchTestData(create bool) db.TableDefinition {
  db.DB_LOCATION = TEST_DB_LOCATION
  if create {
    createTestData()
  }
  return db.FetchTableDefinition(TEST_DB_TABLE_NAME)
}

func removeTestData() {
  logger.Info("Removing %s", TEST_DB_TABLE_NAME)
  db.DeleteTable(TEST_DB_TABLE_NAME)
}

func afterAll() {
  logger.Info("Removing %s", TEST_DB_LOCATION)
  os.RemoveAll(TEST_DB_LOCATION)
}

func TestCreateDelete() {
  removeTestData()
  file := createTestData()

  assert.ToEqual(file.Name, TEST_DB_TABLE_NAME, "Test table should have the correct name")
  assert.ToEqual(file.CSV.LineCount, int64(0), "Test File should have the correct number of lines")

  for ci, col := range TEST_DB_TABLE_COLUMNS {
    assert.ToEqual(file.CSV.Columns[ci], col, "Test Columns should match the file columns")
  }

  assert.ToEqual(file.Info.PrimaryKey, TEST_DB_PK, "Test Table should have correct Primary Key")

  removeTestData()
}

func TestFileParsing() {
  removeTestData()
  file := fetchTestData(true)

  assert.ToEqual(file.Name, TEST_DB_TABLE_NAME, "Test table should have the correct name")
  assert.ToEqual(file.CSV.LineCount, int64(0), "Test File should have the correct number of lines")

  for ci, col := range TEST_DB_TABLE_COLUMNS {
    assert.ToEqual(file.CSV.Columns[ci], col, "Test Columns should match the file columns")
  }

  assert.ToEqual(file.Info.PrimaryKey, TEST_DB_PK, "Test Table should have correct Primary Key")

  removeTestData()
}

func TestInsertion() {
  removeTestData()
  file := createTestData()

  sampleRow := []interface{}{
    "Hello World",
    1234,
    "Hello Aaron",
  }

  nAdded := file.AddData([][]interface{}{ sampleRow })
  assert.ToEqual(nAdded, 1, "Should add the correct number of lines")

  data := file.RetrieveLine(0)
  assert.ToEqual(data[0], int64(0), "Should have correct line numbers")

  for di, val := range sampleRow {
    assert.ToEqual(data[di + 1], val, "Should have inserted the correct values")
  }

  removeTestData()
}

func displayFile() {
  f, err := os.Open(TEST_DB_LOCATION + "/" + TEST_DB_TABLE_NAME + ".csv")
  if err != nil {
    panic(err)
  }

  var buffer []byte = make([]byte, 1024)
  f.Read(buffer)

  logger.Info(string(buffer))
}

func TestDeletion() {
  removeTestData()
  file := createTestData()

  sampleRows := [][]interface{}{
    {
      "hello world",
      1234,
      "hello aaron",
    },
    {
      "idk something, I guess",
      42069,
      "something different",
    },
  }

  nAdded := file.AddData(sampleRows)
  assert.ToEqual(nAdded, 2, "Should have inserted the correct number of lines")

  data := file.RetrieveLine(0)
  for di, val := range sampleRows[0] {
    assert.ToEqual(data[di + 1], val, "Inserting multiple rows should insert in the right order")
  }

  file.DeleteDataByKey(1234)
  assert.ToEqual(file.CSV.LineCount, int64(1), "Should have deleted the line")

  data = file.RetrieveLine(0)

  for di, val := range sampleRows[1] {
    assert.ToEqual(data[di + 1], val, "Should have cascaded the correct values")
  }

  removeTestData()
}

func TestModification() {
  removeTestData()
  file := createTestData()

  file.AddData([][]interface{}{
    {
      "hello world",
      1234,
      "hello aaron",
    },
  })

  file.ModifyRowColumn(0, "Test2", 42069)

  data := file.RetrieveLine(0)
  assert.ToEqual(data[2].(int), 42069, "Should have updated the correct value")

  removeTestData()
}
