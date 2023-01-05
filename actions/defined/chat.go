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
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			loggedIn := player.ConnLoggedIn(conn)

			if loggedIn {
				terminal.AppendGameMessage(conn, message)
			} else {
				chat.SendSystemMessage(conn, message)
			}

			response := parsing.CommandResponse{
				Conn:       conn,
				Person:     true,
				SaveCursor: true,
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

func CreateLocalChatAction(conn net.Conn, message string) definitions.Action {
	return definitions.Action{
		Name:        "Local Message",
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			loggedIn := player.ConnLoggedIn(conn)
			username := player.GetConnUsername(conn)

			response := parsing.CommandResponse{
				LoggedIn:   loggedIn,
				Conn:       conn,
				Chat:       true,
				Global:     true,
				SaveCursor: true,
			}

			chat.SendGlobalMessage(username, message)

			return response
		},
	}
}

func CreateDirectMessageAction(conn net.Conn, target string, message string) definitions.Action {
	return definitions.Action{
		Name:        "Direct Message",
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			tConn := player.LoggedInPlayerMap[target]

			chat.SendMentionMessage(tConn,
				username, target,
				message)

			return parsing.CommandResponse{
				Chat:     true,
				Person:   true,
				Specific: []net.Conn{tConn},
			}
		},
	}
}
