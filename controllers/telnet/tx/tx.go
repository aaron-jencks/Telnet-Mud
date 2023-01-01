package tx

import (
	"mud/services/chat"
	"mud/services/terminal"
	"net"
	"sync"
)

var Clients []net.Conn
var ClientLock sync.Mutex = sync.Mutex{}

func RegisterConnection(conn net.Conn) int {
	ClientLock.Lock()
	Clients = append(Clients, conn)
	chat.RegisterConnection(conn)
	terminal.RegisterConnection(conn)
	cid := len(Clients) - 1
	ClientLock.Unlock()
	return cid
}

func RemoveClient(conn int) {
	ClientLock.Lock()
	chat.UnregisterConnection(Clients[conn])
	terminal.RegisterConnection(Clients[conn])
	Clients = append(Clients[:conn], Clients[conn:]...)
	ClientLock.Unlock()
}

func SendGlobal(body []byte) {
	ClientLock.Lock()
	for _, client := range Clients {
		client.Write(body)
	}
	ClientLock.Unlock()
}

func SendOthers(body []byte, avoid net.Conn) {
	ClientLock.Lock()
	for _, client := range Clients {
		if client != avoid {
			client.Write(body)
		}
	}
	ClientLock.Unlock()
}

func SendTarget(body []byte, target net.Conn) {
	target.Write(body)
}
