package terminal

import (
	"mud/entities"
	"mud/parsing_services/player"
	"mud/services/room"
	"net"
)

type Terminal struct {
	Room   entities.Room
	Buffer []string
}

var TerminalMap map[net.Conn]*Terminal = make(map[net.Conn]*Terminal)

func RegisterConnection(conn net.Conn) {
	TerminalMap[conn] = &Terminal{}
}

func UnregisterConnection(conn net.Conn) {
	delete(TerminalMap, conn)
}

func LoadPlayer(conn net.Conn, username string) {
	p := player.CRUD.Retrieve(username).(entities.Player)
	rint := room.CRUD.Retrieve(p.Room)

	if rint != nil {
		r := rint.(entities.Room)
		TerminalMap[conn] = &Terminal{
			Room: r,
		}
		EnterRoom(conn, r)
	}
}

func AppendGameMessage(conn net.Conn, m string) {
	TerminalMap[conn].Buffer = append(TerminalMap[conn].Buffer, m)
}

func EnterRoom(conn net.Conn, r entities.Room) {
	AppendGameMessage(conn, r.Description)
}

func ChangeRoom(conn net.Conn, r entities.Room) {
	TerminalMap[conn].Room = r
	EnterRoom(conn, r)
}
