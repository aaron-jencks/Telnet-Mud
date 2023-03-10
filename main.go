package main

import (
	"mud/controllers/telnet/rx"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/handlers"
	"mud/utils/net/telnet"
	"net"
)

func main() {
	entities.SetupTables()

	telnet.OptionMap[34] = func(c net.Conn, b []byte) ([]byte, int) {
		return []byte{}, 0
	}

	parsing.RegisterHandler("login", handlers.HandleLogin)
	parsing.RegisterHandler("logout", handlers.HandleLogout)
	parsing.RegisterHandler("register", handlers.HandleRegister)
	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("room", handlers.RoomCrudHandler)
	parsing.RegisterHandler("item", handlers.ItemCrudHandler)
	parsing.RegisterHandler("loot", handlers.LootCrudHandler)
	parsing.RegisterHandler("command", handlers.CommandCrudHandler)
	parsing.RegisterHandler("note", handlers.NoteCrudHandler)
	parsing.RegisterHandler("inventory", handlers.ListInventoryHandler)
	parsing.RegisterHandler("tile", handlers.TileCrudHandler)
	parsing.RegisterHandler("variant", handlers.VariantCrudHandler)
	parsing.RegisterHandler("map", handlers.MapCrudHandler)
	parsing.RegisterHandler("pickup", handlers.HandlePickup)
	parsing.RegisterHandler("about", handlers.HandleInfo)
	parsing.RegisterHandler("here", handlers.ListLootHandler)

	rx.TelnetListenAndServe()
}
