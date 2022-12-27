package command

import (
	"mud/utils/crud"
	"mud/utils/io/db"
	mstrings "mud/utils/strings"
	"regexp"
	"strings"
)

type ExpandedCommand struct {
	Name string
	Args []string
}

func commandToArr(rs map[string]interface{}) []interface{} {
	argsArr := rs["Args"].([]interface{})
	var sargs []string = make([]string, len(argsArr))
	for ai, arg := range argsArr {
		sargs[ai] = arg.(string)
	}
	return []interface{}{
		rs["Name"],
		len(argsArr),
		FormatRegexFromArr(sargs),
	}
}

func commandFromArr(arr []interface{}) interface{} {
	args := FormatRegexToArr(arr[3].(string))
	return ExpandedCommand{
		arr[1].(string),
		args,
	}
}

func createCommandFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	argsArr := args[1].([]string)
	return []interface{}{args[0], len(argsArr), FormatRegexFromArr(argsArr)}
}

var CRUD crud.Crud = crud.CreateCrud("commands", commandToArr, commandFromArr, createCommandFunc)

func CommandExists(name string) bool {
	table := CRUD.FetchTable()
	result := table.Query(name, "Name")
	return len(result) > 0
}

func FormatRegexFromArr(args []string) string {
	formattedStrings := make([]string, len(args))
	for ai, arg := range args {
		formattedStrings[ai] = strings.ReplaceAll(arg, ",", "\\,")
	}
	return strings.Join(formattedStrings, ",")
}

func FormatRegexToArr(argString string) []string {
	var result []string
	var lastStart int = 0
	var isEscaped bool = false

	for bi, b := range argString {
		if isEscaped {
			isEscaped = false
			continue
		}

		switch b {
		case '\\':
			isEscaped = true
		case ',':
			result = append(result, strings.ReplaceAll(argString[lastStart:bi], "\\,", ","))
			lastStart = bi + 1
		}
	}
	// Fetch the last entry
	result = append(result, strings.ReplaceAll(argString[lastStart:], "\\,", ","))

	return result
}

func MatchesCommand(data string, cmd string) bool {
	bits := mstrings.SplitWithQuotes(data, ' ')
	command := CRUD.Retrieve(cmd).(ExpandedCommand)

	if len(bits) != len(command.Args) {
		return false
	} else {
		// Check the regex
		for i := 0; i < len(command.Args); i++ {
			r := regexp.MustCompile(command.Args[i])
			m := r.FindAllStringIndex(bits[i], 1)
			if m == nil || m[0][0] > 0 {
				return false
			}
		}

		return true
	}
}
