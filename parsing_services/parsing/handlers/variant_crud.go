package handlers

import (
	"fmt"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/services/chat"
	"mud/services/variant"
	"mud/utils/strings"
	"net"
)

func HandleVariantCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "variant", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: variant create [id] \"name\" \"icon\""
		if CheckMinArgs(conn, args, 3, usageString) {
			return result
		}

		var nr entities.TileVariant

		if len(args) == 4 {
			idParsed, id := ParseIntegerCheck(conn, args[1], usageString, "id")
			if !idParsed {
				return result
			}

			nr = variant.CRUD.Create(
				id,
				strings.StripQuotes(args[2]),
				parsing.ParseIconString(strings.StripQuotes(args[3]))).(entities.TileVariant)
		} else {
			nr = variant.CRUD.Create(
				strings.StripQuotes(args[1]),
				parsing.ParseIconString(strings.StripQuotes(args[2]))).(entities.TileVariant)
		}

		chat.SendSystemMessage(conn, fmt.Sprintf("Variant %d(%s) created!", nr.Id, nr.Name))

	case "retrieve":
		if CheckMinArgs(conn, args, 3, "Usage: variant retrieve id \"name\"") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: variant retrieve id \"name\"", "id")
		if !idParsed {
			return result
		}

		r := variant.GetSpecificVariant(id, strings.StripQuotes(args[2]))
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Variant:\nId: %d\nName: \"%s\"\nIcon: \"%s\"",
				r.Id, r.Name, r.Icon))

	case "update":
		chat.SendSystemMessage(conn, "Updating a tile variant is not currently supported, please delete and replace")

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: variant delete id \"name\"") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: variant retrieve id \"name\"", "id")
		if !idParsed {
			return result
		}

		name := strings.StripQuotes(args[2])

		variant.CRUD.Delete(id, "Id", name, "Name")
		chat.SendSystemMessage(conn, fmt.Sprintf("Variant %d(%s) deleted!", id, name))
	}

	return result
}
