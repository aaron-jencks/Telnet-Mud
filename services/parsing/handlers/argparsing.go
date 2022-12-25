package handlers

import (
	"mud/services/chat"
	"net"
)

func CheckArgs(conn net.Conn, args []string, target int, message string) bool {
	if len(args) != target {
		chat.SendSystemMessage(conn, message)
		return true
	}
	return false
}

func CheckMinArgs(conn net.Conn, args []string, target int, message string) bool {
	if len(args) < target {
		chat.SendSystemMessage(conn, message)
		return true
	}
	return false
}
