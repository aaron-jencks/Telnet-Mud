package controllers

import (
	"mud/entities"
	"mud/services/chat"
	"mud/services/command"
	"mud/services/parsing"
	"mud/services/player"
	"mud/services/room"
	"mud/services/terminal"
	"mud/services/transition"
	"mud/utils/strings"
	"mud/utils/ui/logger"
	"net"
)

func HandlePacket(conn net.Conn, data []byte) parsing.CommandResponse {
	logger.Info("Parsing: %v", string(data))

	var bits []string = strings.SplitWithQuotes(string(data), ' ')

	if len(bits) == 0 {
		return parsing.CommandResponse{}
	}

	handler, ok := parsing.CommandMap[bits[0]]

	if !(ok || command.CommandExists(bits[0])) {
		chat.SendSystemMessage(conn, "Unknown command, please check your spelling or try again...")
		return parsing.CommandResponse{
			Person: true,
		}
	}

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
				p.Room = newRoom.Id
				player.CRUD.Update(p.Name, p)
			} else {
				chat.SendSystemMessage(conn, "Nothing happened.")
			}
		} else {
			chat.SendSystemMessage(conn, "You need to login first!")
		}

		return parsing.CommandResponse{
			Person: true,
		}
	}
}
