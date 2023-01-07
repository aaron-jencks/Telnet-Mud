package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/utils"
	"net"
)

func CreateScreenBlip(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:        "Blip",
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			return utils.GetDefaultCommandResponse(conn)
		},
	}
}

func CreateGlobalMapRepaint(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:        "Map Repaint",
		AlwaysValid: true,
		Handler: func() parsing.CommandResponse {
			response := utils.GetDefaultMapCommandResponse(conn)
			response.Global = true
			return response
		},
	}
}
