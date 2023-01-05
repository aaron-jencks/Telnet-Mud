package handlers

import (
	"mud/actions/defined"
	"mud/parsing_services/player"
	"net"
)

func ListInventoryHandler(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)
	player.PushAction(username, defined.CreateInventoryListAction(conn))
}

func ListLootHandler(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)
	player.PushAction(username, defined.CreateListLootAction(conn))
}
