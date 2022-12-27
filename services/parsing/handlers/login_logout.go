package handlers

import (
	"mud/services/chat"
	"mud/services/parsing"
	"mud/services/player"
	"mud/services/terminal"
	"net"
)

func HandleLogin(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{
		Person: true,
	}

	if player.ConnLoggedIn(conn) {
		chat.SendSystemMessage(conn, "You're already logged in, you can't log in again!")
	}

	if CheckArgs(conn, args, 2, "Usage: login username password") {
		return result
	}

	if !player.LoginPlayer(args[0], args[1], conn) {
		chat.SendSystemMessage(conn, "Sorry, either that account doesn't exist or the password is incorrect")
	} else {
		terminal.LoadPlayer(conn, args[0])
		chat.SendDirectMessage(conn, "SYSTEM", "Welcome! Please be respectful.")
	}

	return result
}

func HandleLogout(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{}

	if CheckArgs(conn, args, 0, "Usage: logout") {
		return result
	}

	if !player.LogoutPlayer(player.PlayerConnectionMap[conn]) {
		chat.SendSystemMessage(conn, "Sorry, either that account doesn't exist or isn't currently logged in")
		result.Person = true
	}

	return result
}

func HandleRegister(conn net.Conn, args []string) parsing.CommandResponse {
	var result parsing.CommandResponse = parsing.CommandResponse{}

	if player.ConnLoggedIn(conn) {
		chat.SendSystemMessage(conn, "You're already logged in, you can't register a new user from here!")
	}

	if CheckArgs(conn, args, 2, "Usage: register username password") {
		return result
	}

	result.Person = true

	if !player.RegisterPlayer(args[0], args[1]) {
		chat.SendSystemMessage(conn, "Sorry, that account already exists")
	} else {
		chat.SendSystemMessage(conn, "User created successfully, you may now login.")
	}

	return result
}
