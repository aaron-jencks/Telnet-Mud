package controllers

import (
	"mud/services/player"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"mud/utils/ui/pages/chat"
	"net"
	"strings"
)

func GetDisplayForConn(conn net.Conn, saveCursor bool) string {
	var result string = gui.ResetCursorPosition()

	if !player.ConnLoggedIn(conn) {
		result += MOTD()
	}

	result += chat.GetConnChatWindow(conn)

	if saveCursor {
		result = ui.SaveAndResetCursor(result)
	}

	return strings.ReplaceAll(result, "\n", "\n\r")
}

func MOTD() string {
	return "Welcome! Please login using the 'login' command or\n\rcreate a new account using the 'register' command.\n\r"
}
