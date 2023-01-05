package handlers

import (
	"fmt"
	acrud "mud/actions/defined/crud"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/command"
	"mud/utils/strings"
	"net"
)

var CommandCrudHandler parsing.CommandHandler = acrud.CreateCrudParser(
	"command",
	"Usage: command create name \"arg0regex,arg1regex,...,argNregex\"",
	"Usage: command retrieve name",
	"Usage: command update name property:(name|args) \"newValue\"",
	"Usage: command delete name",
	3, 2, 4, 2,
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool { return true },
	func(c net.Conn, s []string) bool { return true },
	func(s []string) []interface{} { return []interface{}{s[0], strings.StripQuotes(s[1])} },
	func(s []string) interface{} { return s[0] },
	func(i interface{}) string { return fmt.Sprintf("Command %s created!", i.(entities.Command).Name) },
	func(i interface{}) string {
		c := i.(command.ExpandedCommand)
		return fmt.Sprintf("Command %s:\nArgs: %v", c.Name, c.Args)
	},
	func(i interface{}) string { return fmt.Sprintf("Command %s updated!", i.(entities.Command).Name) },
	func(i interface{}) string { return fmt.Sprintf("Command %s deleted!", i.(entities.Command).Name) },
	[]string{"name", "args"}, 2,
	func(i interface{}, s1 string, s2 []string) interface{} {
		c := i.(command.ExpandedCommand)

		nv := strings.StripQuotes(s2[0])
		switch s1 {
		case "name":
			c.Name = nv
		case "args":
			c.Args = command.FormatRegexToArr(nv)
		}

		return c
	},
	acrud.DefaultCrudModes, command.CRUD,
)
