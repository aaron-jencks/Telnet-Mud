package handlers

import (
	"mud/entities"
	"mud/services/chat"
	"mud/services/detail"
	"mud/services/parsing"
	"mud/services/player"
	"mud/services/room"
	"mud/utils"
	mstrings "mud/utils/strings"
	"net"
	"strings"
)

func HandleLook(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	var direction string = ""
	var parsedDirections []string
	for _, adb := range args {
		parsedDirections = append(parsedDirections, mstrings.StripQuotes(adb))
	}
	direction = strings.Join(parsedDirections, " ")

	p := player.CRUD.Retrieve(player.PlayerConnectionMap[conn]).(entities.Player)

	details := detail.GetRoomDetails(room.CRUD.Retrieve(p.Room).(entities.Room))
	for _, cDetail := range details {
		if cDetail.Direction == direction {
			if detail.TestPerception(cDetail.Perception, utils.CHECK_DIE) {
				chat.SendSystemMessage(conn, cDetail.Detail)
				return result
			}

			chat.SendSystemMessage(conn, utils.LOOK_FAIL_MESSAGE)
			return result
		}
	}

	chat.SendSystemMessage(conn, "There is nothing in that direction.")
	return result
}
