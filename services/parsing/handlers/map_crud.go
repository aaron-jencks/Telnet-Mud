package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/tmap"
	"mud/utils/strings"
	"net"
)

func HandleMapCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "map", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: map create room \"icon\" x y [z]"

		if CheckMinArgs(conn, args, 5, usageString) {
			return result
		}

		idParsed, rid := ParseIntegerCheck(conn, args[1], usageString, "room")
		if !idParsed {
			return result
		}

		xParsed, x := ParseIntegerCheck(conn, args[3], usageString, "x")
		if !xParsed {
			return result
		}

		yParsed, y := ParseIntegerCheck(conn, args[4], usageString, "y")
		if !yParsed {
			return result
		}

		var nr entities.Map
		if len(args) == 6 {
			zParsed, z := ParseIntegerCheck(conn, args[5], usageString, "z")
			if !zParsed {
				return result
			}

			nr = tmap.CRUD.Create(rid, strings.StripQuotes(args[2]), x, y, z).(entities.Map)
		} else {
			nr = tmap.CRUD.Create(rid, strings.StripQuotes(args[2]), x, y).(entities.Map)
		}

		chat.SendSystemMessage(conn, fmt.Sprintf("Tile %s placed at (Room: %d, X: %d, Y: %d, Z: %d) created!",
			nr.Tile, nr.Room, nr.X, nr.Y, nr.Z))

	case "retrieve":
		usageString := "Usage: map retrieve room x y [z]"

		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		idParsed, rid := ParseIntegerCheck(conn, args[1], usageString, "room")
		if !idParsed {
			return result
		}

		xParsed, x := ParseIntegerCheck(conn, args[2], usageString, "x")
		if !xParsed {
			return result
		}

		yParsed, y := ParseIntegerCheck(conn, args[3], usageString, "y")
		if !yParsed {
			return result
		}

		var r entities.Map

		if len(args) == 5 {
			zParsed, z := ParseIntegerCheck(conn, args[4], usageString, "z")
			if !zParsed {
				return result
			}

			r = tmap.GetTileForCoord(rid, x, y, z)
		} else {
			r = tmap.GetTopMostTile(rid, x, y)
		}

		chat.SendSystemMessage(conn,
			fmt.Sprintf("Map:\nRoom: %d\nCoord: (%d, %d, %d)\nTile: \"%s\"",
				r.Room, r.X, r.Y, r.Z, r.Tile))

	case "update":
		chat.SendSystemMessage(conn, "Updating a tile is not currently supported, please delete and replace")

	case "delete":
		usageString := "Usage: map delete room x y z"
		if CheckMinArgs(conn, args, 5, usageString) {
			return result
		}

		idParsed, rid := ParseIntegerCheck(conn, args[1], usageString, "room")
		if !idParsed {
			return result
		}

		xParsed, x := ParseIntegerCheck(conn, args[2], usageString, "x")
		if !xParsed {
			return result
		}

		yParsed, y := ParseIntegerCheck(conn, args[3], usageString, "y")
		if !yParsed {
			return result
		}

		zParsed, z := ParseIntegerCheck(conn, args[4], usageString, "z")
		if !zParsed {
			return result
		}

		tmap.CRUD.Delete(rid, "Room", x, "X", y, "Y", z, "Z")
		chat.SendSystemMessage(conn, fmt.Sprintf("Map (%d, %d, %d, %d) deleted!", rid, x, y, z))
	}

	return result
}
