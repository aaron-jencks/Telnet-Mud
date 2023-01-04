package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/chat"
	"mud/services/terminal"
	"net"
)

func CreateInfoAction(conn net.Conn, message string) definitions.Action {
	return definitions.Action{
		Name:        "System Message",
		Duration:    0,
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			loggedIn := player.ConnLoggedIn(conn)

			if loggedIn {
				terminal.AppendGameMessage(conn, message)
			} else {
				chat.SendSystemMessage(conn, message)
			}

			response := parsing.CommandResponse{
				Conn:   conn,
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
