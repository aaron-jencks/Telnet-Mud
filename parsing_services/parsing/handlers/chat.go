package handlers

import (
	"mud/actions/defined"
	"mud/parsing_services/player"
	"net"
	"strings"
)

func HandleChat(conn net.Conn, body []string) {
	if len(body) == 0 {
		return
	}

	username := player.GetConnUsername(conn)

	if body[0][0] == '@' {
		// Direct message
		if !player.PlayerLoggedIn(body[0][1:]) {
			player.EnqueueAction(username,
				defined.CreateInfoAction(conn, "That player doese not exist, or is not online"))
		}

		player.EnqueueAction(username,
			defined.CreateDirectMessageAction(conn, body[0][1:], strings.Join(body[1:], " ")))
	} else {
		// Local chat
		player.EnqueueAction(username,
			defined.CreateLocalChatAction(conn, strings.Join(body, " ")))
	}

	player.EnqueueAction(username, defined.CreateScreenBlip(conn))
}
