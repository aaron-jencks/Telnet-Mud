package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/room"
	"mud/services/terminal"
	"mud/utils/strings"
	"net"
)

func HandleRoomCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	if CrudChecks(conn, "room", args) {
		return result
	}

	switch args[0] {
	case "create":
		if CheckMinArgs(conn, args, 3, "Usage: room create \"name\" \"description\"") {
			return result
		}

		nr := room.CRUD.Create(
			strings.StripQuotes(args[1]),
			strings.StripQuotes(args[2])).(entities.Room)
		chat.SendSystemMessage(conn, fmt.Sprintf("Room %d(%s) created!", nr.Id, nr.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: room retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: room retrieve id", "id")
		if !idParsed {
			return result
		}

		r := room.CRUD.Retrieve(id).(entities.Room)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Room %d:\nName: \"%s\"\nDescription: \"%s\"",
				r.Id, r.Name, r.Description))

	case "update":
		if CheckMinArgs(conn, args, 4, "Usage: room update id (name|description) \"newValue\"") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: room update id (name|description) \"newValue\"", "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"name", "description"},
			"Usage: room update id property \"newValue\"", "property") {
			return result
		}

		r := room.CRUD.Retrieve(id).(entities.Room)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "name":
			r.Name = nv
		case "description":
			r.Description = nv
		}

		nr := room.CRUD.Update(id, r).(entities.Room)
		chat.SendSystemMessage(conn, fmt.Sprintf("Room %d(%s) updated!", nr.Id, nr.Name))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: room delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: room delete id", "id")
		if !idParsed {
			return result
		}

		room.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Room %d deleted!", id))
	}

	return result
}

func HandleInfo(conn net.Conn, args []string) parsing.CommandResponse {
	t := terminal.TerminalMap[conn]
	chat.SendSystemMessage(conn, fmt.Sprintf("%s:\n%s", t.Room.Name, t.Room.Description))
	return parsing.CommandResponse{
		Person: true,
	}
}
