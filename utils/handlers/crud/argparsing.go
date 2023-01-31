package crud

import (
	"fmt"
	"mud/actions/defined"
	"mud/parsing_services/player"
	mstrings "mud/utils/strings"
	"net"
	"strings"
)

func CheckArgs(conn net.Conn, args []string, target int, message string) bool {
	username := player.GetConnUsername(conn)

	if len(args) != target {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, message))
		return true
	}
	return false
}

func CheckMinArgs(conn net.Conn, args []string, target int, message string) bool {
	username := player.GetConnUsername(conn)

	if len(args) < target {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, message))
		return false
	}
	return true
}

func RequiresLoggedIn(conn net.Conn) bool {
	username := player.GetConnUsername(conn)

	if !player.ConnLoggedIn(conn) {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, "You must be logged in to perform that action"))
		return true
	}
	return false
}

func CrudChecks(conn net.Conn, crudName string, args []string) bool {
	if RequiresLoggedIn(conn) ||
		!CheckMinArgs(conn, args, 1,
			fmt.Sprintf("Usage: %s (create|retrieve|update|delete) ...", crudName)) ||
		!CheckStringOptions(conn, args[0], []string{"create", "retrieve", "update", "delete"},
			fmt.Sprintf("Usage: %s operation ...", crudName), "operation") {
		return true
	}
	return false
}

func ParseIntegerCheck(conn net.Conn, s string, usageString string, paramName string) (bool, int) {
	username := player.GetConnUsername(conn)

	var id int
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("%s (%s is an integer)", usageString, paramName)))
		return false, -1
	}
	return true, id
}

func CheckStringOptions(conn net.Conn, s string, options []string, usageString, paramName string) bool {
	username := player.GetConnUsername(conn)

	if !mstrings.StringContains(options, s) {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("%s (%s is one of (%s)", usageString, paramName, strings.Join(options, "|"))))
		return false
	}
	return true
}

func ParseBooleanCheck(conn net.Conn, s string, usageString string, paramName string) (bool, bool) {
	username := player.GetConnUsername(conn)

	var id bool
	_, err := fmt.Sscanf(s, "%t", &id)
	if err != nil {
		player.EnqueueAction(username, defined.CreateInfoAction(conn, fmt.Sprintf("%s (%s is an integer)", usageString, paramName)))
		return false, false
	}
	return true, id
}
