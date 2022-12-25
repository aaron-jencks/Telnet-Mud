package parsing

import (
	"mud/services/chat"
	"mud/utils/ui/logger"
	"net"
	"strings"
)

type CommandResponse struct {
	Person   bool
	Global   bool
	Specific []string
}

type DirectMessageMap map[string]string

type CommandHandler func(net.Conn, []string) CommandResponse

var CommandMap map[string]CommandHandler = map[string]CommandHandler{}

func HandlePacket(conn net.Conn, data []byte) CommandResponse {
	logger.Info("Parsing: %v", data)

	var bits []string = strings.Split(string(data), " ")

	if len(bits) == 0 {
		return CommandResponse{}
	}

	handler, ok := CommandMap[bits[0]]

	if ok {
		return handler(conn, bits[1:])
	} else {
		chat.SendSystemMessage(conn, "Unknown command, please check your spelling or try again...")
		return CommandResponse{
			Person: true,
		}
	}
}

func RegisterHandler(command string, handler CommandHandler) {
	logger.Info("Registering %s Command", command)
	CommandMap[command] = handler
}
