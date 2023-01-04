package handlers

import (
	"fmt"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/chat"
	"mud/services/inventory"
	"mud/services/loot"
	"mud/services/room"
	"mud/utils/strings"
	"net"
)

func HandlePickup(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Chat:   true,
		Person: true,
	}

	if CheckMinArgs(conn, args, 1, "Usage: pickup \"item name\" [qty=1]") {
		return result
	}

	p := player.CRUD.Retrieve(player.PlayerConnectionMap[conn]).(entities.Player)
	r := room.CRUD.Retrieve(p.Room).(entities.Room)
	roomLoot := loot.GetLootForRoom(r)

	var qty int = 1
	if len(args) > 1 {
		idParsed, pQty := ParseIntegerCheck(conn, args[1], "Usage: pickup \"item name\" [qty=1]", "qty")
		if !idParsed {
			return result
		}
		qty = pQty
	}

	for _, loot := range roomLoot {
		if loot.Item.Name == strings.StripQuotes(args[0]) && loot.Quantity <= qty {
			nQty := inventory.AddItemToInventory(p, loot.Item, qty)
			chat.SendSystemMessage(conn, fmt.Sprintf("You now have %dx %s", nQty, loot.Item.Name))
			return result
		} else if loot.Quantity < qty {
			chat.SendSystemMessage(conn, fmt.Sprintf("There are only %dx here", loot.Quantity))
			return result
		}
	}

	chat.SendSystemMessage(conn, "There is none of that here to pick up")
	return result
}
