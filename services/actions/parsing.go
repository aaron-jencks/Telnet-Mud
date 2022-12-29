package actions

import (
	"fmt"
	"mud/entities"
	"mud/services/actions/handlers"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
	"net"
)

var letterMapping map[rune]func(entities.Player) = map[rune]func(entities.Player){
	'u': handlers.HandleUpMovement,
	'l': handlers.HandleLeftMovement,
	'r': handlers.HandleRightMovement,
	'd': handlers.HandleDownMovement,
}

func ParseString(conn net.Conn, s string) parsing.CommandResponse {
	p := player.CRUD.Retrieve(player.PlayerConnectionMap[conn]).(entities.Player)

	for _, r := range s {
		handler, ok := letterMapping[r]
		if !ok {
			chat.SendSystemMessage(conn, fmt.Sprintf("Unknown parsing symbol %c", r))
			return parsing.CommandResponse{
				Chat:   true,
				Person: true,
			}
		} else {
			handler(p)
		}
	}

	return parsing.CommandResponse{}
}
