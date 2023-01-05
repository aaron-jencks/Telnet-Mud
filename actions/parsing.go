package actions

import (
	"fmt"
	"mud/actions/handlers"
	"mud/parsing_services/parsing"
	"mud/services/chat"
	"net"
)

var letterMapping map[rune]func(net.Conn) = map[rune]func(net.Conn){
	'u': handlers.HandleUpMovement,
	'l': handlers.HandleLeftMovement,
	'r': handlers.HandleRightMovement,
	'd': handlers.HandleDownMovement,
}

func ParseString(conn net.Conn, s string) parsing.CommandResponse {
	var num string = ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			num += string(r)
			continue
		} else {
			handler, ok := letterMapping[r]
			if !ok {
				chat.SendSystemMessage(conn, fmt.Sprintf("Unknown parsing symbol %c", r))
				return parsing.CommandResponse{
					Chat:   true,
					Person: true,
				}
			} else {
				var inum int
				fmt.Sscanf(num, "%d", &inum)

				for i := 0; i < inum; i++ {
					handler(conn)
				}
				num = ""
			}
		}
	}

	return parsing.CommandResponse{}
}
