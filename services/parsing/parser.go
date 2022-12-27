package parsing

import (
	"mud/entities"
	"mud/services/chat"
	"mud/services/command"
	"mud/services/player"
	"mud/services/room"
	"mud/services/terminal"
	"mud/services/transition"
	"mud/utils/strings"
	"mud/utils/ui/logger"
	"net"
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
	logger.Info("Parsing: %v", string(data))

	var bits []string = strings.SplitWithQuotes(string(data), ' ')

	if len(bits) == 0 {
		return CommandResponse{}
	}

	if !command.CommandExists(bits[0]) {
		chat.SendSystemMessage(conn, "Unknown command, please check your spelling or try again...")
		return CommandResponse{
			Person: true,
		}
	}

	handler, ok := CommandMap[bits[0]]

	if ok {
		return handler(conn, bits[1:])
	} else {
		if player.ConnLoggedIn(conn) {
			username := player.PlayerConnectionMap[conn]
			p := player.CRUD.Retrieve(username).(entities.Player)
			exists, trans := transition.TransitionExists(p.Room, bits[0], bits[1:])
			if exists {
				newRoom := room.CRUD.Retrieve(trans.Target).(entities.Room)
				terminal.ChangeRoom(conn, newRoom)
			} else {
				chat.SendSystemMessage(conn, "Nothing happened.")
			}
		} else {
			chat.SendSystemMessage(conn, "You need to login first!")
		}

		return CommandResponse{
			Person: true,
		}
	}
}

func RegisterHandler(command string, handler CommandHandler) {
	logger.Info("Registering %s Command", command)
	CommandMap[command] = handler
}
