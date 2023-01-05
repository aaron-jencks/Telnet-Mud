package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"net"
)

func CreateScreenBlip(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:        "Blip",
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			return parsing.CommandResponse{
				Person:   true,
				Conn:     conn,
				LoggedIn: player.ConnLoggedIn(conn),
			}
		},
	}
}
