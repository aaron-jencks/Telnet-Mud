package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/inventory"
	"mud/services/parsing"
	"mud/services/player"
	"net"
	"strings"
)

func ListInventoryHandler(conn net.Conn, args []string) parsing.CommandResponse {
	items := inventory.GetPlayerInventory(player.CRUD.Retrieve(player.PlayerConnectionMap[conn]).(entities.Player))

	var displayList []string
	for _, item := range items {
		displayList = append(displayList, fmt.Sprintf("%dx %s", item.Quantity, item.Item.Name))
	}

	text := "Inventory: Empty."
	if len(displayList) > 0 {
		text = "Inventory:\n" + strings.Join(displayList, "\n")
	}
	chat.SendSystemMessage(conn, text)

	return parsing.CommandResponse{
		Person: true,
	}
}
