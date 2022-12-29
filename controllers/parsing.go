package controllers

import (
	"mud/services/actions"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
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

	if ok {
		return handler(conn, bits[1:])
	} else {
		if player.ConnLoggedIn(conn) {
			return actions.ParseString(conn, string(data))
		} else {
			chat.SendSystemMessage(conn, "You need to login first!")
			return parsing.CommandResponse{
				Chat:   true,
				Person: true,
			}
		}
	}
}
