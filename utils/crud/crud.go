package crud

import (
	"errors"
	"fmt"
	"mud/utils/io/db"
)

func CreateCrud(tableName string, selectorFormatter SelectorFormatter, toArrFunc ToArrFunc, scannerFunc db.RowScanner, rowSelector RowSelector, fromArrFunc FromArrFunc, createFunc CreateFunc) Crud {
	return Crud{tableName, selectorFormatter, toArrFunc, scannerFunc, rowSelector, fromArrFunc, createFunc}
}

func (c Crud) FetchTable() db.TableDefinition {
	_, ok := tableMap[c.TableName]
	if !ok {
		newTable := db.FetchTableDefinition(c.TableName)
		tableMap[c.TableName] = newTable
	}
	return tableMap[c.TableName]
}

func (c Crud) Create(args ...interface{}) interface{} {
	table := c.FetchTable()
	newValue := c.createFunction(table, args...)
	_, rc := table.AddData([][]interface{}{
		newValue,
	})

	if rc > 1 || rc == 0 {
		panic(errors.New("Create function for CRUD didn't insert a new row, or inserted too many"))
	}

	return c.Retrieve(c.rowSelector(newValue)...)
}

func (c Crud) Retrieve(args ...interface{}) interface{} {
	table := c.FetchTable()
	query := fmt.Sprintf("select from %s where %s", c.TableName, c.selectorFormatter(args))
	results := table.QueryData(query, c.scannerFunc)
	if len(results) > 0 {
		result := results[0]
		return result
	}
	return nil
}

func (c Crud) RetrieveAll(args ...interface{}) []interface{} {
	table := c.FetchTable()
	query := fmt.Sprintf("select from %s where %s", c.TableName, c.selectorFormatter(args))
	results := table.QueryData(query, c.scannerFunc)
	if len(results) > 0 {
		return results
	}
	return nil
}

func (c Crud) Update(newValue interface{}, selectors ...interface{}) interface{} {
	table := c.FetchTable()
	oldRows := table.QueryData(fmt.Sprintf("select from %s where %s", c.TableName, c.selectorFormatter(selectors)), c.scannerFunc)
	if len(oldRows) > 0 {
		line := oldRows[0][0].(int)
		table.ModifyRow(line, c.toArrFunc(newValue))
		return newValue
	}
	return nil
}

func (c Crud) UpdateQuery(retrieveValues []interface{}, retrieveColumns []string, newValue interface{}) interface{} {
	table := c.FetchTable()

	var queryArgs []interface{} = make([]interface{}, len(retrieveColumns)<<1)

	for ai := range retrieveValues {
		argsIndex := ai << 1
		queryArgs[argsIndex] = retrieveValues[ai]
		queryArgs[argsIndex+1] = retrieveColumns[ai]
	}

	oldRows := table.MultiQuery(queryArgs...)
	if len(oldRows) > 0 {
		line := oldRows[0][0].(int)
		table.ModifyRow(line, c.toArrFunc(newValue))
		return newValue
	}

	return nil
}

func (c Crud) Delete(keys ...interface{}) {
	table := c.FetchTable()

	if len(keys) == 1 {
		table.DeleteDataByKey(keys[0])
	} else {
		lines := table.MultiQuery(keys...)

		var lineNumbers []int64 = make([]int64, len(lines))

		for li, line := range lines {
			lineNumbers[li] = int64(line[0].(int))
		}

		table.DeleteLines(lineNumbers)
	}
}
