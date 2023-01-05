package defined

import (
	"fmt"
	"mud/actions/definitions"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/inventory"
	"mud/services/loot"
	"mud/services/room"
	"net"
	"strings"
	"time"
)

func CreateInventoryListAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "List Inventory",
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			items := inventory.GetPlayerInventory(player.CRUD.Retrieve(username).(entities.Player))

			var displayList []string
			for _, item := range items {
				displayList = append(displayList, fmt.Sprintf("%dx %s", item.Quantity, item.Item.Name))
			}

			text := "Inventory: Empty."
			if len(displayList) > 0 {
				text = "Inventory:\n" + strings.Join(displayList, "\n")
			}

			player.PushAction(username, CreateInfoAction(conn, text))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
		},
	}
}

func CreatePickupItemAction(conn net.Conn, p entities.Player, targetItem string, qty int) definitions.Action {
	return definitions.Action{
		Name:       "Pickup",
		Duration:   1 * time.Second,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			var result parsing.CommandResponse = parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}

			r := room.CRUD.Retrieve(p.Room).(entities.Room)
			roomLoot := loot.GetLootForRoom(r)

			for _, loot := range roomLoot {
				if loot.Item.Name == targetItem && loot.Quantity <= qty {
					nQty := inventory.AddItemToInventory(p, loot.Item, qty)
					player.PushAction(p.Name, CreateInfoAction(conn, fmt.Sprintf("You now have %dx %s", nQty, loot.Item.Name)))
					return result
				} else if loot.Quantity < qty {
					player.PushAction(p.Name, CreateInfoAction(conn, fmt.Sprintf("There are only %dx here", loot.Quantity)))
					return result
				}
			}

			player.PushAction(p.Name, CreateInfoAction(conn, "There is none of that here to pick up"))
			return result
		},
	}
}

func CreateListLootAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Pickup",
		Duration:   500 * time.Millisecond,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			p := player.CRUD.Retrieve(username).(entities.Player)
			r := room.CRUD.Retrieve(p.Room).(entities.Room)
			roomLoot := loot.GetLootForPosition(r, p.RoomX, p.RoomY)

			var displayList []string
			for _, item := range roomLoot {
				displayList = append(displayList, fmt.Sprintf("%dx %s", item.Quantity, item.Item.Name))
			}

			text := "Loot: Empty."
			if len(displayList) > 0 {
				text = "Loot:\n" + strings.Join(displayList, "\n")
			}

			player.PushAction(username, CreateInfoAction(conn, text))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
		},
	}
}
