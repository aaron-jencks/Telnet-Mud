package actions

import (
	"fmt"
	"mud/actions/handlers"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"mud/services/chat"
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
