package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/loot"
	"mud/services/parsing"
	"net"
)

func HandleLootCrud(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CrudChecks(conn, "loot", args) {
		return result
	}

	switch args[0] {
	case "create":
		usageString := "Usage: loot create room item quantity"
		if CheckMinArgs(conn, args, 4, usageString) {
			return result
		}

		idParsed, rId := ParseIntegerCheck(conn, args[1], usageString, "room")
		if !idParsed {
			return result
		}

		idParsed, iId := ParseIntegerCheck(conn, args[2], usageString, "item")
		if !idParsed {
			return result
		}

		qtyParsed, qty := ParseIntegerCheck(conn, args[3], usageString, "quantity")
		if !qtyParsed {
			return result
		}

		nr := loot.CRUD.Create(rId, iId, qty).(entities.Loot)
		chat.SendSystemMessage(conn, fmt.Sprintf("Loot %d created!", nr.Id))

	case "retrieve":
		if CheckMinArgs(conn, args, 2, "Usage: loot retrieve id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: loot retrieve id", "id")
		if !idParsed {
			return result
		}

		r := loot.CRUD.Retrieve(id).(entities.Loot)
		chat.SendSystemMessage(conn,
			fmt.Sprintf("Loot %d:\nRoom: %d\nItem: %d\nQuantity: %d",
				r.Id, r.Room, r.Item, r.Quantity))

	case "update":
		if CheckMinArgs(conn, args, 4, "Usage: loot update id (room|item|quantity) newValue") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: loot update id (room|item|quantity) newValue", "id")
		if !idParsed {
			return result
		}

		if CheckStringOptions(conn, args[2], []string{"room", "item", "quantity"},
			"Usage: loot update id property newValue", "property") {
			return result
		}

		idParsed, newValue := ParseIntegerCheck(conn, args[3], "Usage: loot update id (room|item|quantity) newValue", "newValue")
		if !idParsed {
			return result
		}

		r := loot.CRUD.Retrieve(id).(entities.Loot)
		switch args[2] {
		case "room":
			r.Room = newValue
		case "item":
			r.Item = newValue
		case "quantity":
			r.Quantity = newValue
		}

		nr := loot.CRUD.Update(id, r).(entities.Loot)
		chat.SendSystemMessage(conn, fmt.Sprintf("Loot %d updated!", nr.Id))

	case "delete":
		if CheckMinArgs(conn, args, 2, "Usage: loot delete id") {
			return result
		}

		idParsed, id := ParseIntegerCheck(conn, args[1], "Usage: loot delete id", "id")
		if !idParsed {
			return result
		}

		loot.CRUD.Delete(id)
		chat.SendSystemMessage(conn, fmt.Sprintf("Loot %d deleted!", id))
	}

	return result
}
