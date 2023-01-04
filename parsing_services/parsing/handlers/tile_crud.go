package handlers

import (
	"fmt"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/chat"
	"mud/services/tile"
	"mud/utils/strings"
	"net"
)

func HandleTileCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "tile", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: tile create \"name\" \"type\" \"icon\" [bg fg]"
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		name := strings.StripQuotes(args[1])
		itype := strings.StripQuotes(args[2])
		icon := parsing.ParseIconString(strings.StripQuotes(args[3]))
		var nr entities.Tile

		if len(args) == 6 {
			bgParsed, bg := ParseIntegerCheck(conn, args[4], usageString, "bg")
			if !bgParsed {
				return result
			}

			fgParsed, fg := ParseIntegerCheck(conn, args[5], usageString, "fg")
			if !fgParsed {
				return result
			}

			nr = tile.CRUD.Create(name, itype, icon, bg, fg).(entities.Tile)
		} else {
			nr = tile.CRUD.Create(name, itype, icon).(entities.Tile)
		}

		chat.SendSystemMessage(conn, fmt.Sprintf("Tile %s created!", nr.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: tile retrieve \"name\"") {
			return result
		}

		r := tile.CRUD.Retrieve(strings.StripQuotes(args[1])).(entities.Tile)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Tile:\nName: \"%s\"\nType: \"%s\"\nIcon: \"\033[%dm\033%dm%s\033[0m\"",
				r.Name, r.IconType, r.BG, r.FG, r.Icon))

	case "update":
		usageString := "Usage: tile update \"name\" (name|type|icon|bg|fg) \"newValue\""
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"name", "type", "icon", "bg", "fg"},
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
		case "bg":
			bgParsed, bg := ParseIntegerCheck(conn, nv, usageString, "bg")
			if !bgParsed {
				return result
			}

			r.BG = bg
		case "fg":
			fgParsed, fg := ParseIntegerCheck(conn, nv, usageString, "fg")
			if !fgParsed {
				return result
			}

			r.FG = fg
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
