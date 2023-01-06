package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/utils"
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

			response := utils.GetDefaultCommandResponse(conn)
			response.SaveCursor = true

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
			username := player.GetConnUsername(conn)

			response := utils.GetDefaultChatCommandResponse(conn)

			response.Person = false
			response.Global = true
			response.SaveCursor = true

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

			return utils.GetDefaultChatCommandResponse(conn)
		},
	}
}
