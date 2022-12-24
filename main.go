package main

import (
	"mud/entities"
	"mud/services/parsing"
	"mud/services/parsing/handlers"
	"mud/utils/net/telnet"
	"mud/utils/ui/logger"
)

func main() {
	entities.SetupTables()

	logger.Info("Starting Server on port 8080")

	parsing.RegisterHandler("chat", handlers.HandleChat)
	parsing.RegisterHandler("global", handlers.HandleGlobal)

	telnet.TelnetListenAndServe()
}
