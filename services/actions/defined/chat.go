package defined

import (
	"mud/entities"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
	"mud/utils/actions"
	"net"
)

func CreateInfoAction(conn net.Conn, message string) actions.Action {
	return actions.Action{
		Name:        "System Message",
		Duration:    0,
		AlwaysValid: true,
		Handler: func(p entities.Player) parsing.CommandResponse {
			chat.SendSystemMessage(conn, message)

			response := parsing.CommandResponse{
				Person: true,
			}

			if player.ConnLoggedIn(conn) {
				response.Info = true
			} else {
				response.Chat = true
			}

			return response
		},
	}
}
