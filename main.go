package main

import (
	"mud/controllers/telnet/rx"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/handlers"
)

func main() {
	entities.SetupTables()

	parsing.RegisterHandler("login", handlers.HandleLogin)
	parsing.RegisterHandler("logout", handlers.HandleLogout)
	parsing.RegisterHandler("register", handlers.HandleRegister)
	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("room", handlers.RoomCrudHandler)
	parsing.RegisterHandler("item", handlers.ItemCrudHandler)
	parsing.RegisterHandler("command", handlers.CommandCrudHandler)
	parsing.RegisterHandler("note", handlers.NoteCrudHandler)
	// parsing.RegisterHandler("inventory", handlers.ListInventoryHandler)
	// parsing.RegisterHandler("tile", handlers.HandleTileCrud)
	parsing.RegisterHandler("map", handlers.MapCrudHandler)
	// parsing.RegisterHandler("pickup", handlers.HandlePickup)
	parsing.RegisterHandler("about", handlers.HandleInfo)

	rx.TelnetListenAndServe()
}
