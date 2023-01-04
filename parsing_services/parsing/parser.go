package parsing

import (
	"mud/utils/ui/logger"
	"net"
)

type CommandResponse struct {
	Map      bool
	Chat     bool
	Info     bool
	Clear    bool
	Person   bool
	Global   bool
	Specific []string
}

type DirectMessageMap map[string]string

type CommandHandler func(net.Conn, []string)

var CommandMap map[string]CommandHandler = map[string]CommandHandler{}

func RegisterHandler(command string, handler CommandHandler) {
	logger.Info("Registering %s Command", command)
	CommandMap[command] = handler
}
