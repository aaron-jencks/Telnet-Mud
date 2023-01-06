package actions

import (
	"fmt"
	"mud/actions/defined"
	"mud/actions/handlers"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"net"
)

var letterMapping map[rune]func(net.Conn) = map[rune]func(net.Conn){
	'u': handlers.HandleUpMovement,
	'l': handlers.HandleLeftMovement,
	'r': handlers.HandleRightMovement,
	'd': handlers.HandleDownMovement,
}

func ParseString(conn net.Conn, s string) parsing.CommandResponse {
	username := player.GetConnUsername(conn)

	var num string = ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			num += string(r)
			continue
		} else {
			handler, ok := letterMapping[r]
			if !ok {
				player.PushAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("Unknown parsing symbol %c", r)))
				return parsing.CommandResponse{
					Conn:   conn,
					Chat:   true,
					Person: true,
				}
			} else {
				var inum int = 1

				if len(num) > 0 {
					fmt.Sscanf(num, "%d", &inum)
				}

				for i := 0; i < inum; i++ {
					handler(conn)
				}

				player.PushAction(username, defined.CreateScreenBlip(conn))
				num = ""
			}
		}
	}

	return parsing.CommandResponse{}
}
