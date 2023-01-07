package handlers

import (
	"mud/actions/defined"
	"mud/entities"
	"mud/parsing_services/player"
	"mud/utils/handlers/crud"
	"mud/utils/strings"
	"net"
)

func HandlePickup(conn net.Conn, args []string) {
	if crud.CheckMinArgs(conn, args, 1, "Usage: pickup \"item name\" [qty=1]") {
		return
	}

	username := player.GetConnUsername(conn)
	p := player.CRUD.Retrieve(username).(entities.Player)

	var qty int = 1
	if len(args) > 1 {
		idParsed, pQty := crud.ParseIntegerCheck(conn, args[1], "Usage: pickup \"item name\" [qty=1]", "qty")
		if !idParsed {
			return
		}
		qty = pQty
	}

	targetItem := strings.StripQuotes(args[0])

	player.EnqueueAction(username, defined.CreatePickupItemAction(conn, p, targetItem, qty))
}
