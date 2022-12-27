package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/item"
	"mud/services/parsing"
	"mud/utils/strings"
	"net"
)

func HandleItemCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	if CrudChecks(conn, "item", args) {
		return result
	}

	switch args[0] {
	case "create":
		if CheckMinArgs(conn, args, 3, "Usage: item create \"name\" \"description\"") {
			return result
		}

		nr := item.CRUD.Create(
			strings.StripQuotes(args[1]),
			strings.StripQuotes(args[2])).(entities.Item)
		chat.SendSystemMessage(conn, fmt.Sprintf("Item %d(%s) created!", nr.Id, nr.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: item retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: item retrieve id", "id")
		if !idParsed {
			return result
		}

		r := item.CRUD.Retrieve(id).(entities.Item)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Item %d:\nName: \"%s\"\nDescription: \"%s\"",
				r.Id, r.Name, r.Description))

	case "update":
		if CheckMinArgs(conn, args, 4, "Usage: item update id (name|description) \"newValue\"") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: item update id (name|description) \"newValue\"", "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"name", "description"},
			"Usage: item update id property \"newValue\"", "property") {
			return result
		}

		r := item.CRUD.Retrieve(id).(entities.Item)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "name":
			r.Name = nv
		case "description":
			r.Description = nv
		}

		nr := item.CRUD.Update(id, r).(entities.Item)
		chat.SendSystemMessage(conn, fmt.Sprintf("Item %d(%s) updated!", nr.Id, nr.Name))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: item delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: item delete id", "id")
		if !idParsed {
			return result
		}

		item.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Item %d deleted!", id))
	}

	return result
}
