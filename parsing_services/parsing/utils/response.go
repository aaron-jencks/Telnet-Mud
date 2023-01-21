package utils

import (
	"mud/parsing_services/parsing"
	"mud/parsing_services/player"
	"net"
)

func GetDefaultCommandResponse(conn net.Conn) parsing.CommandResponse {
	username := player.GetConnUsername(conn)
	loggedIn := player.ConnLoggedIn(conn)

	result := parsing.CommandResponse{
		LoggedIn: loggedIn,
		Conn:     conn,
		Person:   true,
	}

	if loggedIn {
		result.Player = player.FetchPlayerByName(username)
	}

	return result
}

func GetDefaultInfoCommandResponse(conn net.Conn) parsing.CommandResponse {
	result := GetDefaultCommandResponse(conn)
	result.Info = true
	return result
}

func GetDefaultChatCommandResponse(conn net.Conn) parsing.CommandResponse {
	result := GetDefaultCommandResponse(conn)
	result.Chat = true
	return result
}

func GetDefaultMapCommandResponse(conn net.Conn) parsing.CommandResponse {
	result := GetDefaultCommandResponse(conn)
	result.Map = true
	return result
}

func GetDefaultRepaintCommandResponse(conn net.Conn) parsing.CommandResponse {
	result := GetDefaultInfoCommandResponse(conn)
	result.Chat = true
	result.Map = true
	return result
}
