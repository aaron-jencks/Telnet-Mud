package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/tile"
	"mud/utils/strings"
	"net"
)

func HandleTileCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	if CrudChecks(conn, "tile", args) {
		return result
	}

	switch args[0] {
	case "create":
		if CheckMinArgs(conn, args, 4, "Usage: tile create \"name\" \"type\" \"icon\"") {
			return result
		}

		nr := tile.CRUD.Create(
			strings.StripQuotes(args[1]),
			strings.StripQuotes(args[2]),
			parsing.ParseIconString(strings.StripQuotes(args[3]))).(entities.Tile)
		chat.SendSystemMessage(conn, fmt.Sprintf("Tile %s created!", nr.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: tile retrieve \"name\"") {
			return result
		}

		r := tile.CRUD.Retrieve(strings.StripQuotes(args[1])).(entities.Tile)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Tile:\nName: \"%s\"\nType: \"%s\"\nIcon: \"%s\"",
				r.Name, r.IconType, r.Icon))

	case "update":
		if CheckMinArgs(conn, args, 4, "Usage: tile update \"name\" (name|type|icon) \"newValue\"") {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"name", "type", "icon"},
			"Usage: tile update \"name\" property \"newValue\"", "property") {
			return result
		}

		id := strings.StripQuotes(args[1])

		r := tile.CRUD.Retrieve(id).(entities.Tile)
		nv := strings.StripQuotes(args[3])
		switch args[2] {
		case "name":
			r.Name = nv
		case "type":
			r.IconType = nv
		case "icon":
			r.Icon = parsing.ParseIconString(nv)
		}

		nr := tile.CRUD.Update(id, r).(entities.Tile)
		chat.SendSystemMessage(conn, fmt.Sprintf("Tile %s updated!", nr.Name))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: tile delete \"name\"") {
			return result
		}

		id := strings.StripQuotes(args[1])

		tile.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Tile %s deleted!", id))
	}

	return result
}
