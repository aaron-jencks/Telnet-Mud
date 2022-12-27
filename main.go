package main

import (
	"mud/entities"
	"mud/services/parsing"
	"mud/services/parsing/handlers"
	"mud/utils/net/telnet"
)

func main() {
	entities.SetupTables()

	parsing.RegisterHandler("login", handlers.HandleLogin)
	parsing.RegisterHandler("logout", handlers.HandleLogout)
	parsing.RegisterHandler("register", handlers.HandleRegister)
	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("room", handlers.HandleRoomCrud)

	telnet.TelnetListenAndServe()
}
