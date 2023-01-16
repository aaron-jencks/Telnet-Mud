package command

import (
	"database/sql"
	"fmt"
	"mud/entities"
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

func commandToArr(rs interface{}) []interface{} {
	rec := rs.(ExpandedCommand)
	return []interface{}{
		rec.Name,
		len(rec.Args),
		FormatRegexFromArr(rec.Args),
	}
}

func commandFromArr(arr []interface{}) interface{} {
	args := FormatRegexToArr(arr[3].(string))
	return ExpandedCommand{
		arr[1].(string),
		args,
	}
}

func createCommandFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	argsArr := args[1].([]string)
	return []interface{}{args[0], len(argsArr), FormatRegexFromArr(argsArr)}
}

func commandSelector(row []interface{}) string {
	return fmt.Sprintf("Name=\"%s\"", row[0].(string))
}

func commandScanner(row *sql.Rows) (interface{}, error) {
	var result entities.Command = entities.Command{}
	err := row.Scan(&result.Name, &result.ArgCount, &result.ArgRegex)
	if err != nil {
		return nil, err
	}
	return result, err
}

func commandUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ocs := oldValue.(entities.Command)
	ncs := newValue.(entities.Command)

	var result []crud.RowModStruct
	if ocs.Name != ncs.Name {
		result = append(result, crud.RowModStruct{
			Column:   "Name",
			NewValue: fmt.Sprintf("\"%s\"", ncs.Name),
		})
	}
	if ocs.ArgCount != ncs.ArgCount {
		result = append(result, crud.RowModStruct{
			Column:   "ArgCount",
			NewValue: ncs.ArgCount,
		})
	}
	if ocs.ArgRegex != ncs.ArgRegex {
		result = append(result, crud.RowModStruct{
			Column:   "ArgRegex",
			NewValue: fmt.Sprintf("\"%s\"", ncs.ArgRegex),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("commands",
	commandSelector, commandToArr, commandScanner, commandFromArr,
	createCommandFunc, commandUpdateFunc)

func CommandExists(name string) bool {
	table := CRUD.FetchTable()
	result := table.QueryData(fmt.Sprintf("Name=\"%s\"", name), commandScanner)
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
