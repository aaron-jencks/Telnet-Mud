package crud

import (
	"mud/utils/io/db"
)

var tableMap map[string]*db.TableDefinition = make(map[string]*db.TableDefinition)

type ICrud interface {
	Create(...interface{}) interface{}
	Retrieve(interface{}) interface{}
	Update(interface{}, interface{}) interface{}
	Delete(interface{})
}

type CrudUpdate struct {
	Key     interface{}
	NewData interface{}
}

type Crud struct {
	TableName      string
	toArrFunc      func(interface{}) []interface{}
	fromArrFunc    func([]interface{}) interface{}
	createFunction func(*db.TableDefinition, ...interface{}) []interface{}
}

func CreateCrud(tableName string, toArrFunc func(interface{}) []interface{}, fromArrFunc func([]interface{}) interface{}, createFunc func(*db.TableDefinition, ...interface{}) []interface{}) Crud {
	return Crud{tableName, toArrFunc, fromArrFunc, createFunc}
}

func (c Crud) FetchTable() *db.TableDefinition {
	_, ok := tableMap[c.TableName]
	if !ok {
		newTable := db.FetchTableDefinition(c.TableName)
		tableMap[c.TableName] = &newTable
	}
	return tableMap[c.TableName]
}

func (c Crud) Create(args ...interface{}) interface{} {
	table := c.FetchTable()
	newValue := c.createFunction(table, args...)
	table.AddData([][]interface{}{
		newValue,
	})
	return c.Retrieve(newValue[table.Info.PrimaryKey])
}

func (c Crud) Retrieve(value interface{}) interface{} {
	table := c.FetchTable()
	results := table.QueryPK(value)
	if len(results) > 0 {
		result := results[0]
		return c.fromArrFunc(result)
	}
	return nil
}

func (c Crud) RetrieveAll(value interface{}) []interface{} {
	table := c.FetchTable()
	results := table.QueryPK(value)
	if len(results) > 0 {
		var result []interface{} = make([]interface{}, len(results))
		for ri, row := range results {
			result[ri] = c.fromArrFunc(row)
		}
		return result
	}
	return nil
}

func (c Crud) Update(retrieveValue interface{}, newValue interface{}) interface{} {
	table := c.FetchTable()
	oldRows := table.QueryPK(retrieveValue)
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
