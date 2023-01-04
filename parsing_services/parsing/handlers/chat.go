package handlers

import (
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	chatService "mud/services/chat"
	"net"
	"strings"
)

func HandleChat(conn net.Conn, body []string) parsing.CommandResponse {
	if len(body) == 0 {
		return parsing.CommandResponse{}
	}

	result := parsing.CommandResponse{
		Chat: true,
	}

	var senderName string
	if !player.ConnLoggedIn(conn) {
		senderName = "Anonymous"
	} else {
		senderName = player.PlayerConnectionMap[conn]
	}

	if body[0][0] == '@' {
		// Direct message
		if !player.PlayerLoggedIn(body[0][1:]) {
			result.Person = true
			chatService.SendSystemMessage(conn, "That player doese not exist, or is not online")
		}

		tConn := player.LoggedInPlayerMap[body[0][1:]]

		chatService.SendMentionMessage(tConn,
			senderName, body[0],
			strings.Join(body[1:], " "))

		result.Specific = append(result.Specific, tConn)
	} else {
		// Local chat
		chatService.SendGlobalMessage(senderName, strings.Join(body, " "))
		result.Global = true
	}

	return result
}
