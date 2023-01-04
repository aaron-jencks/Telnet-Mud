package handlers

import (
	"fmt"
	"mud/parsing_services/parsing"
	"mud/services/chat"
	"mud/services/command"
	"mud/utils/strings"
	"net"
)

func HandleCommandCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "command", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: command create name \"arg0regex,arg1regex,...,argNregex\""
		if CheckMinArgs(conn, args, 3, usageString) {
			return result
		}

		nc := command.CRUD.Create(args[1], strings.StripQuotes(args[2])).(command.ExpandedCommand)
		chat.SendSystemMessage(conn, fmt.Sprintf("Command %s created!", nc.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: command retrieve name") {
			return result
		}

		c := command.CRUD.Retrieve(args[1]).(command.ExpandedCommand)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Command %s:\nArgs: %v", c.Name, c.Args))

	case "update":
		usageString := "Usage: command update name (name|args) \"newValue\""
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"name", "args"},
			"Usage: command update name property \"newValue\"", "property") {
			return result
		}

		c := command.CRUD.Retrieve(args[1]).(command.ExpandedCommand)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "name":
			c.Name = nv
		case "args":
			c.Args = command.FormatRegexToArr(nv)
		}

		nc := command.CRUD.Update(args[1], c).(command.ExpandedCommand)
		chat.SendSystemMessage(conn, fmt.Sprintf("Command %s updated!", nc.Name))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: command delete name") {
			return result
		}

		command.CRUD.Delete(args[1])
		chat.SendSystemMessage(conn, fmt.Sprintf("Command %s deleted!", args[1]))
	}

	return result
}
