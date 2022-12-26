package controllers

import (
	"mud/services/player"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"mud/utils/ui/pages/chat"
	"net"
	"strings"
)

func GetDisplayForConn(conn net.Conn, saveCursor, clearScreen bool) string {
	var result string

	if clearScreen {
		result += gui.Clearscreen()
	} else {
		result += gui.ResetCursorPosition()
	}

	yStart := 0
	if !player.ConnLoggedIn(conn) {
		result += MOTD()
		yStart = 2
	}

	result += gui.AnsiOffsetText(40, yStart, chat.GetConnChatWindow(conn))

	if saveCursor {
		result = ui.SaveAndResetCursor(result)
	}

	return strings.ReplaceAll(result, "\n", "\n\r")
}

func MOTD() string {
	return strings.Join(
		ui.CreateTextParagraph(
			"Welcome! Please login using the 'login' command or create a new account using the 'register' command.\n\r",
			utils.WINDOW_W),
		"\n\r")
}
