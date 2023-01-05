package controllers

import (
	"mud/actions"
	"mud/actions/defined"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/utils/strings"
	"mud/utils/ui/logger"
	"net"
)

func HandlePacket(conn net.Conn, data []byte) {
	logger.Info("Parsing: %v", string(data))

	var bits []string = strings.SplitWithQuotes(string(data), ' ')

	if len(bits) == 0 {
		return
	}

	handler, ok := parsing.CommandMap[bits[0]]

	username := player.GetConnUsername(conn)

	if ok {
		handler(conn, bits[1:])
		player.PushAction(username,
			defined.CreateScreenBlip(conn))
	} else {
		if player.ConnLoggedIn(conn) {
			actions.ParseString(conn, string(data))
		} else {
			player.PushAction(player.GetAnonymousUsername(conn), defined.CreateInfoAction(conn, "You need to login first!"))
		}
	}
}
