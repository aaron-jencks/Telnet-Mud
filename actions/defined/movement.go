package defined

import (
	"mud/actions/definitions"
	"mud/entities"
	"mud/parsing_services/parsing"
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
			p := player.CRUD.Retrieve(username).(entities.Player)

			if p.RoomY > 0 {
				p.RoomY--
				player.CRUD.Update(username, p)

				return parsing.CommandResponse{
					LoggedIn: true,
					Map:      true,
					Person:   true,
				}
			}

			player.PushAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
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
			p := player.CRUD.Retrieve(username).(entities.Player)

			if p.RoomX > 0 {
				p.RoomX--
				player.CRUD.Update(username, p)

				return parsing.CommandResponse{
					LoggedIn: true,
					Map:      true,
					Person:   true,
				}
			}

			player.PushAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
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
			p := player.CRUD.Retrieve(username).(entities.Player)
			r := room.CRUD.Retrieve(p.Room).(entities.Room)

			if p.RoomX < r.Width-1 {
				p.RoomX++
				player.CRUD.Update(username, p)

				return parsing.CommandResponse{
					LoggedIn: true,
					Map:      true,
					Person:   true,
				}
			}

			player.PushAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
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
			p := player.CRUD.Retrieve(username).(entities.Player)
			r := room.CRUD.Retrieve(p.Room).(entities.Room)

			if p.RoomY < r.Height-1 {
				p.RoomY++
				player.CRUD.Update(username, p)

				return parsing.CommandResponse{
					LoggedIn: true,
					Map:      true,
					Person:   true,
				}
			}

			player.PushAction(username, CreateInfoAction(conn, "You're at the edge of the room"))

			return parsing.CommandResponse{
				LoggedIn: true,
				Info:     true,
				Person:   true,
			}
		},
	}
}
