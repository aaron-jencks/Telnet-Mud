package defined

import (
	"mud/entities"
	"mud/services/parsing"
	"mud/utils/actions"
	"net"
)

func SystemMessageHandler(conn net.Conn, message string) actions.ActionHandler {
	return func(p entities.Player) parsing.CommandResponse {

	}
}

func CreateInfoAction(conn net.Conn, message string) actions.Action {
	return actions.Action{
		Name:        "System Message",
		Duration:    0,
		AlwaysValid: true,
	}
}
