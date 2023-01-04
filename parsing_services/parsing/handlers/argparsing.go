package handlers

import (
	"fmt"
	"mud/parsing_services/player"
	"mud/services/chat"
	mstrings "mud/utils/strings"
	"net"
	"strings"
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

func RequiresLoggedIn(conn net.Conn) bool {
	if !player.ConnLoggedIn(conn) {
		chat.SendSystemMessage(conn, "You must be logged in to perform that action")
		return true
	}
	return false
}

func CrudChecks(conn net.Conn, crudName string, args []string) bool {
	if RequiresLoggedIn(conn) ||
		CheckMinArgs(conn, args, 1,
			fmt.Sprintf("Usage: %s (create|retrieve|update|delete) ...", crudName)) ||
		CheckStringOptions(conn, args[0], []string{"create", "retrieve", "update", "delete"},
			fmt.Sprintf("Usage: %s operation ...", crudName), "operation") {
		return true
	}
	return false
}

func ParseIntegerCheck(conn net.Conn, s string, usageString string, paramName string) (bool, int) {
	var id int
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		chat.SendSystemMessage(conn, fmt.Sprintf("%s (%s is an integer)", usageString, paramName))
		return false, -1
	}
	return true, id
}

func CheckStringOptions(conn net.Conn, s string, options []string, usageString, paramName string) bool {
	if !mstrings.StringContains(options, s) {
		chat.SendSystemMessage(conn, fmt.Sprintf("%s (%s is one of (%s)", usageString, paramName, strings.Join(options, "|")))
		return true
	}
	return false
}
