package csv_tests

import (
  "logger"
  "test_defs"
  "assert"
  "csv"
  "os"
)

var TEST_CSV_PATH string = "./test.csv"
var TEST_CSV_COLUMNS []string = []string {
  "Test1",
  "Test2",
  "Test3",
}
var TEST_CSV_LINES [][]string = [][]string {
  {
    "\"idk\\sseomthing\"",
    "\"I\\shave\\scommas,,,,,\"",
    "OneWord",
  },
  {
    "1234",
    "Hello\\sworld",
    "something\\selse",
  },
}
var TEST_CSV_LINE_LOCATIONS []int64 = []int64 {
  20,
  68,
  102,
}
var EXPECTED_CSV_FILE_TEXT string = "2\nTest1,Test2,Test3\n\"idk\\sseomthing\",\"I\\shave\\scommas,,,,,\",OneWord\n1234,Hello\\sworld,something\\selse\n"

var TEST_SUITE test_defs.TestSuite = test_defs.TestSuite {
  "CSV File",
  []func() {
    TestFileCreation,
    TestCreateDelete,
    TestFileParsing,
    TestLineContents,
    TestLineModification,
    TestLineDeletion,
  },
}

func createTestData() csv.CSVFile {
  return csv.CreateCSV(
    TEST_CSV_PATH,
    TEST_CSV_COLUMNS,
    TEST_CSV_LINES)
}

func removeTestFile() {
  logger.Info("Removing %s", TEST_CSV_PATH)
  err := os.Remove(TEST_CSV_PATH)
  if err != nil {
    logger.Warn(err)
  }
}

func parseTestFile(create bool) csv.CSVFile {
  if create {
    createTestData()
  }
  return csv.ParseCSV(TEST_CSV_PATH)
}

func TestFileCreation() {
  removeTestFile()
  createTestData()

  f, err := os.Open(TEST_CSV_PATH)
  defer f.Close()
  assert.ToEqual(err, nil, "Creating the CSV file should create the corresponding file")

  expectedLength := len(EXPECTED_CSV_FILE_TEXT)
  buff := make([]byte, expectedLength + 1024)
  nOut, err := f.Read(buff)
  buff = buff[:nOut]

  assert.ToEqual(err, nil, "Program should have read/write permissions for the csv file")
  assert.ToEqual(nOut, expectedLength, "The file contents should be the expected length")

  assert.ToEqual(string(buff), EXPECTED_CSV_FILE_TEXT, "File should be written correctly")

  removeTestFile()
}

func TestCreateDelete() {
  removeTestFile()
  file := createTestData()

  assert.ToEqual(file.LineCount, int64(len(TEST_CSV_LINES)), "Test File should have the correct number of lines")

  for ci, col := range TEST_CSV_COLUMNS {
    assert.ToEqual(file.Columns[ci], col, "Test Columns should match the file columns")
  }

  for li, loc := range TEST_CSV_LINE_LOCATIONS {
    assert.ToEqual(file.LineLocations[li], loc, "Test File should have correct line locations")
  }

  assert.ToEqual(file.Filepath, TEST_CSV_PATH, "Test File should have correct path")

  removeTestFile()
}

func TestFileParsing() {
  removeTestFile()
  file := parseTestFile(true)

  assert.ToEqual(file.LineCount, int64(len(TEST_CSV_LINES)), "Test File should have the correct number of lines")

  for ci, col := range TEST_CSV_COLUMNS {
    assert.ToEqual(file.Columns[ci], col, "Test Columns should match the file columns")
  }

  for li, loc := range TEST_CSV_LINE_LOCATIONS {
    assert.ToEqual(file.LineLocations[li], loc, "Test File should have correct line locations")
  }

  assert.ToEqual(file.Filepath, TEST_CSV_PATH, "Test File should have correct path")

  removeTestFile()
}

func TestLineContents() {
  removeTestFile()
  file := parseTestFile(true)

  for li, line := range TEST_CSV_LINES {
    data := file.ReadSpecificLine(int64(li))
    for ci, col := range line {
      assert.ToEqual(data[ci], col, "Data should be read from the file correctly")
    }
  }

  removeTestFile()
}

func TestLineModification() {
  removeTestFile()
  file := createTestData()
  testData := []string{
    "Hello",
    "This is",
    "New",
  }

  file.ModifyLine(0, testData)

  line := file.ReadSpecificLine(0)
  logger.Info(line)
  for ci, col := range testData {
    assert.ToEqual(line[ci], col, "Modified lines should show up in the file")
  }

  removeTestFile()
}

func TestLineDeletion() {
  removeTestFile()
  file := createTestData()
  file.DeleteLine(0)

  assert.ToEqual(file.LineCount, int64(len(TEST_CSV_LINES) - 1), "LineCount should be updated")
  line := file.ReadSpecificLine(0)
  for ci, col := range TEST_CSV_LINES[1] {
    assert.ToEqual(line[ci], col, "Deleted lines should cause the rest of the file to move up")
  }

  newFile := parseTestFile(false)
  assert.ToEqual(newFile.LineCount, int64(len(TEST_CSV_LINES) - 1), "LineCount should be updated in the file")

  removeTestFile()
}
