package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/detail"
	"mud/services/parsing"
	"mud/utils/strings"
	"net"
)

func HandleDetailCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	if CrudChecks(conn, "detail", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: detail create room \"direction\" \"detail\" perception"
		if CheckMinArgs(conn, args, 5, usageString) {
			return result
		}

		idParsed, rId := ParseIntegerCheck(conn, args[1], usageString, "room")
		if !idParsed {
			return result
		}

		idParsed, pReq := ParseIntegerCheck(conn, args[4], usageString, "perception")
		if !idParsed {
			return result
		}

		nr := detail.CRUD.Create(rId, args[2], strings.StripQuotes(args[3]), pReq).(entities.Detail)
		chat.SendSystemMessage(conn, fmt.Sprintf("Detail %d(%s) created!", nr.Id, nr.Direction))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: detail retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: detail retrieve id", "id")
		if !idParsed {
			return result
		}

		r := detail.CRUD.Retrieve(id).(entities.Detail)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Detail %d:\nRoom: %d\nDirection: \"%s\"\nDetail: \"%s\"\nPerception: %d",
				r.Id, r.Room, r.Direction, r.Detail, r.Perception))

	case "update":
		usageString := "Usage: detail update id (room|direction|detail|perception) \"newValue\""
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], usageString, "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"room", "direction", "detail", "perception"},
			"Usage: detail update id property \"newValue\"", "property") {
			return result
		}

		r := detail.CRUD.Retrieve(id).(entities.Detail)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "room":
			idParsed, rId := ParseIntegerCheck(conn, nv, usageString, "newValue")
			if !idParsed {
				return result
			}
			r.Room = rId
		case "direction":
			r.Direction = nv
		case "detail":
			r.Detail = nv
		case "perception":
			idParsed, pReq := ParseIntegerCheck(conn, nv, usageString, "newValue")
			if !idParsed {
				return result
			}
			r.Perception = pReq
		}

		nr := detail.CRUD.Update(id, r).(entities.Detail)
		chat.SendSystemMessage(conn, fmt.Sprintf("Detail %d(%s) updated!", nr.Id, nr.Direction))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: detail delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: detail delete id", "id")
		if !idParsed {
			return result
		}

		detail.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Detail %d deleted!", id))
	}

	return result
}
