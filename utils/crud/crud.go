package crud

import (
	"errors"
	"fmt"
	"mud/utils/io/db"
	"mud/utils/ui/logger"
	"reflect"
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

func (c Crud) Update(newValue interface{}, selectorArgs ...interface{}) interface{} {
	table := c.FetchTable()
	selector := c.selectorFormatter(selectorArgs)

	oldValue := c.Retrieve(selectorArgs...)
	if oldValue != nil {
		// Convert to reflect container
		oldValue := reflect.ValueOf(&oldValue)
		newValue := reflect.ValueOf(&newValue)
		oldType := reflect.TypeOf(oldValue)
		newType := reflect.TypeOf(newValue)

		if oldValue.NumField() == newValue.NumField() {
			// Extract the struct itself
			oldStruct := oldValue.Elem()
			newStruct := newValue.Elem()

			for fi := range make([]int, oldValue.NumField()) {
				oldField := oldStruct.Field(fi)
				newField := newStruct.Field(fi)

				// for struct field name
				oldFieldType := oldType.Field(fi)
				newFieldType := newType.Field(fi)

				if oldField.Type().Name() == newField.Type().Name() &&
					oldFieldType.Name == newFieldType.Name {
					oldField.Set(newField)

					// Now that the field is updated correctly
					// we can update it in the database
					table.UpdateData(selector,
						oldFieldType.Name,
						newField.Interface())
				} else {
					logger.Error("Modifying column type/name is not supported in update statements")
					return nil
				}
			}
		} else {
			logger.Error("Adding or Removing columns from a table is not supported for Update actions")
			return nil
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
