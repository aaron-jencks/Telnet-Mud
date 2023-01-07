package handlers

import (
	"mud/actions/defined"
	"mud/parsing_services/player"
	"mud/services/chat"
	"mud/utils/handlers/crud"
	"net"
)

func HandleLogin(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if player.ConnLoggedIn(conn) {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, "You're already logged in, you can't log in again!"))
	}

	if crud.CheckArgs(conn, args, 2, "Usage: login username password") {
		return
	}

	player.EnqueueAction(username, defined.CreateLoginAction(conn, args[0], args[1]))
	player.EnqueueAction(username, defined.CreateScreenBlip(conn))
}

func HandleLogout(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if crud.CheckArgs(conn, args, 0, "Usage: logout") {
		return
	}

	player.EnqueueAction(username, defined.CreateLogoutAction(conn))
	player.EnqueueAction(username, defined.CreateScreenBlip(conn))
}

func HandleRegister(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if player.ConnLoggedIn(conn) {
		chat.SendSystemMessage(conn, "You're already logged in, you can't register a new user from here!")
	}

	if crud.CheckArgs(conn, args, 2, "Usage: register username password") {
		return
	}

	player.EnqueueAction(username, defined.CreateRegisterAction(conn, args[0], args[1]))
	player.EnqueueAction(username, defined.CreateScreenBlip(conn))
}
