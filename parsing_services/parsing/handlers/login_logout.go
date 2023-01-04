package handlers

import (
	"mud/actions/defined"
	"mud/parsing_services/player"
	"mud/services/chat"
	"net"
)

func HandleLogin(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if player.ConnLoggedIn(conn) {
		player.PushAction(username, defined.CreateInfoAction(conn, "You're already logged in, you can't log in again!"))
	}

	if CheckArgs(conn, args, 2, "Usage: login username password") {
		return
	}

	player.PushAction(username, defined.CreateLoginAction(conn, args[0], args[1]))
}

func HandleLogout(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if CheckArgs(conn, args, 0, "Usage: logout") {
		return
	}

	player.PushAction(username, defined.CreateLogoutAction(conn))
}

func HandleRegister(conn net.Conn, args []string) {
	username := player.GetConnUsername(conn)

	if player.ConnLoggedIn(conn) {
		chat.SendSystemMessage(conn, "You're already logged in, you can't register a new user from here!")
	}

	if CheckArgs(conn, args, 2, "Usage: register username password") {
		return
	}

	player.PushAction(username, defined.CreateRegisterAction(conn, args[0], args[1]))
}
