package handlers

import (
	"mud/actions/defined"
	"mud/parsing_services/player"
	"net"
)

func HandleUpMovement(conn net.Conn) {
	username := player.GetConnUsername(conn)
	player.EnqueueAction(username, defined.CreateMoveUpAction(conn))
}

func HandleLeftMovement(conn net.Conn) {
	username := player.GetConnUsername(conn)
	player.EnqueueAction(username, defined.CreateMoveLeftAction(conn))
}

func HandleRightMovement(conn net.Conn) {
	username := player.GetConnUsername(conn)
	player.EnqueueAction(username, defined.CreateMoveRightAction(conn))
}

func HandleDownMovement(conn net.Conn) {
	username := player.GetConnUsername(conn)
	player.EnqueueAction(username, defined.CreateMoveDownAction(conn))
}
