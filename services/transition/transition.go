package transition

import (
	"mud/services/command"
	"mud/utils/crud"
	"mud/utils/io/db"
)

type ExpandedTransition struct {
	Id          int
	Source      int
	Target      int
	Command     string
	CommandArgs []string
}

func transitionToArr(rs map[string]interface{}) []interface{} {
	argsArr := rs["CommandArgs"].([]interface{})
	var sargs []string = make([]string, len(argsArr))
	for ai, arg := range argsArr {
		sargs[ai] = arg.(string)
	}
	return []interface{}{
		int(rs["Id"].(float64)),
		int(rs["Source"].(float64)),
		int(rs["Target"].(float64)),
		rs["Command"],
		command.FormatRegexFromArr(sargs),
	}
}

func transitionFromArr(arr []interface{}) interface{} {
	return ExpandedTransition{
		arr[1].(int),
		arr[2].(int),
		arr[3].(int),
		arr[4].(string),
		command.FormatRegexToArr(arr[5].(string)),
	}
}

func createTransitionFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	id := 0
	if table.CSV.LineCount > 0 {
		id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
	}

	return []interface{}{id, args[0], args[1], args[2], command.FormatRegexFromArr(args[3].([]string))}
}

var CRUD crud.Crud = crud.CreateCrud("transitions", transitionToArr, transitionFromArr, createTransitionFunc)
