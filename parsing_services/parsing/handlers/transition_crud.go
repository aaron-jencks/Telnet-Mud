package handlers

import (
	"fmt"
	"mud/parsing_services/parsing"
	"mud/services/chat"
	"mud/services/command"
	"mud/services/transition"
	"mud/utils/strings"
	"net"
)

func HandleTransitionCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "transition", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: transition create source target command \"arg0regex,arg1regex,...,argNregex\""
		if CheckMinArgs(conn, args, 5, usageString) {
			return result
		}

		idParsed, sourceId := ParseIntegerCheck(conn, args[1], usageString, "source")
		if !idParsed {
			return result
		}

		idParsed, targetId := ParseIntegerCheck(conn, args[2], usageString, "target")
		if !idParsed {
			return result
		}

		nt := transition.CRUD.Create(sourceId, targetId, args[3], strings.StripQuotes(args[4])).(transition.ExpandedTransition)
		chat.SendSystemMessage(conn, fmt.Sprintf("Transition %d(%d->%d) created!", nt.Id, nt.Source, nt.Target))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: transition retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: transition retrieve id", "id")
		if !idParsed {
			return result
		}

		t := transition.CRUD.Retrieve(id).(transition.ExpandedTransition)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Transition %d:\nSource: %d\nTarget: %d\nCommand: %s\nArgs: %v",
				t.Id, t.Source, t.Target, t.Command, t.CommandArgs))

	case "update":
		usageString := "Usage: transition update id (source|target|command|args) \"newValue\""
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], usageString, "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"source", "target", "command", "args"},
			"Usage: transition update id property \"newValue\"", "property") {
			return result
		}

		t := transition.CRUD.Retrieve(id).(transition.ExpandedTransition)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "source":
			idParsed, nId := ParseIntegerCheck(conn, nv, usageString, "source")
			if !idParsed {
				return result
			}

			t.Source = nId
		case "target":
			idParsed, nId := ParseIntegerCheck(conn, nv, usageString, "target")
			if !idParsed {
				return result
			}

			t.Target = nId
		case "command":
			t.Command = nv
		case "args":
			t.CommandArgs = command.FormatRegexToArr(nv)
		}

		nt := transition.CRUD.Update(id, t).(transition.ExpandedTransition)
		chat.SendSystemMessage(conn, fmt.Sprintf("Room %d(%d->%d) updated!", nt.Id, nt.Source, nt.Target))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: transition delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: transition delete id", "id")
		if !idParsed {
			return result
		}

		transition.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Transition %d deleted!", id))
	}

	return result
}
