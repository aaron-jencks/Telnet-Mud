package crud

import (
	"errors"
	"mud/utils/io/db"
	"mud/utils/ui/logger"
)

func CreateCrud(tableName string, selectorFormatter SelectorFormatter, toArrFunc ToArrFunc,
	scannerFunc db.RowScanner, fromArrFunc FromArrFunc, createFunc CreateFunc, updatefunc UpdateFunc) Crud {
	return Crud{tableName, selectorFormatter, toArrFunc, scannerFunc, fromArrFunc, createFunc, updatefunc}
}

func (c Crud) FetchTable() db.TableDefinition {
	_, ok := tableMap[c.TableName]
	if !ok {
		newTable := db.FetchTableDefinition(c.TableName)
		tableMap[c.TableName] = newTable
	}
	return tableMap[c.TableName]
}

func (c Crud) Create(args ...interface{}) int64 {
	table := c.FetchTable()
	newValue := c.createFunction(table, args...)
	rid, rc := table.AddData([][]interface{}{
		newValue,
	})

	if rc > 1 || rc == 0 {
		panic(errors.New("Create function for CRUD didn't insert a new row, or inserted too many"))
	}

	return rid[0]
}

func (c Crud) Retrieve(args ...interface{}) interface{} {
	table := c.FetchTable()
	results := table.QueryData(c.selectorFormatter(args), c.scannerFunc)
	if len(results) > 0 {
		result := results[0]
		return result
	}
	return nil
}

func (c Crud) RetrieveAll(args ...interface{}) []interface{} {
	table := c.FetchTable()
	results := table.QueryData(c.selectorFormatter(args), c.scannerFunc)
	if len(results) > 0 {
		return results
	}
	return nil
}

func (c Crud) Update(newValue interface{}, selectorArgs ...interface{}) interface{} {
	table := c.FetchTable()
	selector := c.selectorFormatter(selectorArgs)

	oldValue := c.Retrieve(selectorArgs...)
	if oldValue != nil {
		modifValues := c.updateFunc(oldValue, newValue)
		for mvi := range modifValues {
			// Now that the field is updated correctly
			// we can update it in the database
			table.UpdateData(selector,
				modifValues[mvi].Column,
				modifValues[mvi].NewValue)
		}
	} else {
		logger.Error("You can't insert a new value using an update statement")
		return nil
	}

	return c.Retrieve(selectorArgs...)
}

func (c Crud) Delete(keys ...interface{}) {
	table := c.FetchTable()

	table.DeleteData(c.selectorFormatter(keys))
}
