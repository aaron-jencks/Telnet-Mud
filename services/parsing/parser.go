package parsing

import (
	"mud/utils/ui/logger"
	"strings"
)

type CommandResponse struct {
	Person   string
	Others   string
	Global   string
	Specific map[string]string
}

type DirectMessageMap map[string]string

type CommandHandler func([]string) CommandResponse

var CommandMap map[string]CommandHandler = map[string]CommandHandler{}

func HandlePacket(data []byte) CommandResponse {
	var bits []string = strings.Split(string(data), " ")

	if len(bits) == 0 {
		return CommandResponse{}
	}

	handler, ok := CommandMap[bits[0]]

	if ok {
		return handler(bits[1:])
	} else {
		return CommandResponse{
			Person: "Unknown Command, please check your spelling or try again",
		}
	}
}

func RegisterHandler(command string, handler CommandHandler) {
	logger.Info("Registering %s Command", command)
	CommandMap[command] = handler
}
