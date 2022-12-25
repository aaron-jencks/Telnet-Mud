package main

import (
	"mud/entities"
	"mud/services/parsing"
	"mud/services/parsing/handlers"
	"mud/utils/net/telnet"
)

func main() {
	entities.SetupTables()

	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("global", handlers.HandleGlobal)

	telnet.TelnetListenAndServe()
}
