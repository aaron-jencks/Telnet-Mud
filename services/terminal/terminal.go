package terminal

import (
	"mud/entities"
	"mud/services/room"
	"net"
)

type Terminal struct {
	Room   room.ExpandedRoom
	Buffer []string
}

var TerminalMap map[net.Conn]*Terminal = make(map[net.Conn]*Terminal)

func RegisterConnection(conn net.Conn) {
	TerminalMap[conn] = &Terminal{}
}

func UnregisterConnection(conn net.Conn) {
	delete(TerminalMap, conn)
}

func LoadPlayer(conn net.Conn, p entities.Player) {
	rint := room.CRUD.Retrieve(p.Room)

	if rint != nil {
		r := rint.(room.ExpandedRoom)
		TerminalMap[conn] = &Terminal{
			Room: r,
		}
		EnterRoom(conn, r)
	}
}

func AppendGameMessage(conn net.Conn, m string) {
	TerminalMap[conn].Buffer = append(TerminalMap[conn].Buffer, m)
}

func EnterRoom(conn net.Conn, r room.ExpandedRoom) {
	AppendGameMessage(conn, r.Description)
}

func ChangeRoom(conn net.Conn, r room.ExpandedRoom) {
	TerminalMap[conn].Room = r
	EnterRoom(conn, r)
}
