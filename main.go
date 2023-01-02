package main

import (
	"mud/controllers/telnet/rx"
	"mud/entities"
	"mud/services/parsing"
	"mud/services/parsing/handlers"
	"mud/services/player"
)

func main() {
	entities.SetupTables()

	player.CreateDefaultGlobalHandler()

	parsing.RegisterHandler("login", handlers.HandleLogin)
	parsing.RegisterHandler("logout", handlers.HandleLogout)
	parsing.RegisterHandler("register", handlers.HandleRegister)
	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("room", handlers.HandleRoomCrud)
	parsing.RegisterHandler("item", handlers.HandleItemCrud)
	parsing.RegisterHandler("transition", handlers.HandleTransitionCrud)
	parsing.RegisterHandler("command", handlers.HandleCommandCrud)
	parsing.RegisterHandler("note", handlers.HandleNoteCrud)
	parsing.RegisterHandler("inventory", handlers.ListInventoryHandler)
	parsing.RegisterHandler("tile", handlers.HandleTileCrud)
	parsing.RegisterHandler("map", handlers.HandleMapCrud)
	parsing.RegisterHandler("pickup", handlers.HandlePickup)
	parsing.RegisterHandler("about", handlers.HandleInfo)

	rx.TelnetListenAndServe()
}
