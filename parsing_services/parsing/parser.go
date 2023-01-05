package parsing

import (
	"mud/entities"
	"mud/utils/ui/logger"
	"net"
)

type CommandResponse struct {
	LoggedIn   bool
	Map        bool
	Chat       bool
	Info       bool
	Clear      bool
	Person     bool
	Global     bool
	SaveCursor bool
	Conn       net.Conn
	Player     entities.Player
	Specific   []net.Conn
}

type DirectMessageMap map[string]string

type CommandHandler func(net.Conn, []string)

var CommandMap map[string]CommandHandler = map[string]CommandHandler{}

func RegisterHandler(command string, handler CommandHandler) {
	logger.Info("Registering %s Command", command)
	CommandMap[command] = handler
}
