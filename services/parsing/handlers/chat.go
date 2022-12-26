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

		chatService.SendMentionMessage(player.LoggedInPlayerMap[body[0][1:]],
			senderName, body[0],
			strings.Join(body[1:], " "))

		result.Specific = append(result.Specific, body[0][1:])
	} else {
		// Local chat
		chatService.SendGlobalMessage(senderName, strings.Join(body, " "))
		result.Global = true
	}

	return result
}
