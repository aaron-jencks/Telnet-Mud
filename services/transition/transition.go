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

func TransitionExists(source int, command string, args []string) (bool, ExpandedTransition) {
	table := CRUD.FetchTable()
	result := table.Query(source, "Source")
	for _, row := range result {
		if row[4].(string) == command {
			tid := row[0].(int)
			matches, et := MatchesTransition(tid, command, args)
			if matches {
				return true, et
			}
		}
	}
	return false, ExpandedTransition{}
}

func MatchesTransition(tid int, command string, args []string) (bool, ExpandedTransition) {
	t := CRUD.Retrieve(tid).(ExpandedTransition)

	if t.Command == command && len(args) == len(t.CommandArgs) {
		for ai, arg := range t.CommandArgs {
			if arg[0] == '"' && arg[len(arg)-1] == '"' {
				arg = arg[1 : len(arg)-1]
			}

			if t.CommandArgs[ai] != arg {
				return false, ExpandedTransition{}
			}
		}

		return true, t
	}

	return false, ExpandedTransition{}
}
