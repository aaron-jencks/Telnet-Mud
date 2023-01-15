package db

import "testing"

func TestTableCreateDelete(t *testing.T) {
	DeleteTable("TestTable")

	CreateTableIfNotExist(
		"TestTable",
		[]string{
			"TestColumn1",
			"TestColumn2",
			"TestColumn3",
		},
		[]string{
			"text",
			"text",
			"integer",
		},
	)

	DeleteTable("TestTable")
}
