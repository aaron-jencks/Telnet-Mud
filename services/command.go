package command_service

import (
  "db"
  "crud"
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
