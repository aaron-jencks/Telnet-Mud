package terminal

import (
	"mud/entities"
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

}

func ChangeRoom(conn net.Conn, r entities.Room) {
	TerminalMap[conn].Room = r
}
