package handlers

import (
	chatService "mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
	"net"
	"strings"
)

func HandleChat(conn net.Conn, body []string) parsing.CommandResponse {
	if len(body) == 0 {
		return parsing.CommandResponse{}
	}

	result := parsing.CommandResponse{}

	if body[0][0] == '@' {
		// Direct message
		if !player.PlayerLoggedIn(body[0][1:]) {
			result.Person = true
			chatService.SendSystemMessage(conn, "That player doese not exist, or is not online")
		}

		chatService.SendDirectMessage(player.LoggedInPlayerMap[body[0][1:]],
			player.PlayerConnectionMap[conn],
			strings.Join(body[1:], " "))
	} else {
		// Local chat
		chatService.SendGlobalMessage("Anonymous", strings.Join(body, " "))
		result.Global = true
	}

	return result
}
