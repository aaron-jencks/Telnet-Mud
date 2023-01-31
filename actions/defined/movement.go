package defined

import (
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/utils"
	"mud/parsing_services/player"
	"mud/services/room"
	"net"
	"time"
)

func CreateMoveUpAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Move Up",
		Duration:   100 * time.Millisecond,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			p := player.FetchPlayerByName(username)

			if p.RoomY > 0 {
				p.RoomY--
				player.CRUD.Update(p, p.Id)

				return utils.GetDefaultMapCommandResponse(conn)
			}

			player.EnqueueAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return utils.GetDefaultInfoCommandResponse(conn)
		},
	}
}

func CreateMoveLeftAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Move Left",
		Duration:   100 * time.Millisecond,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			p := player.FetchPlayerByName(username)

			if p.RoomX > 0 {
				p.RoomX--
				player.CRUD.Update(p, p.Id)

				return utils.GetDefaultMapCommandResponse(conn)
			}

			player.EnqueueAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return utils.GetDefaultInfoCommandResponse(conn)
		},
	}
}

func CreateMoveRightAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Move Right",
		Duration:   100 * time.Millisecond,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			p := player.FetchPlayerByName(username)
			r := room.CRUD.Retrieve(p.Room).(room.ExpandedRoom)

			if p.RoomX < r.Width-1 {
				p.RoomX++
				player.CRUD.Update(p, p.Id)

				return utils.GetDefaultMapCommandResponse(conn)
			}

			player.EnqueueAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return utils.GetDefaultInfoCommandResponse(conn)
		},
	}
}

func CreateMoveDownAction(conn net.Conn) definitions.Action {
	return definitions.Action{
		Name:       "Move Down",
		Duration:   100 * time.Millisecond,
		ValidModes: []string{"Logged In"},
		Handler: func() parsing.CommandResponse {
			username := player.GetConnUsername(conn)
			p := player.FetchPlayerByName(username)
			r := room.CRUD.Retrieve(p.Room).(room.ExpandedRoom)

			if p.RoomY < r.Height-1 {
				p.RoomY++
				player.CRUD.Update(p, p.Id)

				return utils.GetDefaultMapCommandResponse(conn)
			}

			player.EnqueueAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return utils.GetDefaultInfoCommandResponse(conn)
		},
	}
}
