package controllers

import (
	"fmt"
	"mud/entities"
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
	"mud/utils/actions"
	"net"
)

func ParseActions(conn net.Conn, actions []actions.Action) parsing.CommandResponse {
	result := parsing.CommandResponse{
		Person: true,
	}

	p := player.PlayerConnectionMap[conn]
	pent := player.CRUD.Retrieve(p).(entities.Player)
	currentMode := pent.CurrentMode

	for _, action := range actions {
		if !action.AlwaysValid {
			found := false
			for _, amode := range action.ValidModes {
				if currentMode == amode {
					found = true
					break
				}
			}

			if !found {
				chat.SendSystemMessage(conn, fmt.Sprintf("%s cannot be performed while in %s mode", action.Name, pent.CurrentMode))
				result.Chat = true
				return result
			}
		}

		// The action can be performed

	}

	return result
}
