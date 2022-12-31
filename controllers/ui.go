package controllers

import (
	"mud/services/parsing"
	"mud/services/player"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"mud/utils/ui/pages/chat"
	"mud/utils/ui/pages/terminal"
	"mud/utils/ui/pages/tmap"
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

	if !player.ConnLoggedIn(conn) {
		result += MOTD()
		result += chat.GetConnChatWindowModHeight(conn, utils.CHAT_H-2) + "\n"
	} else {
		result += gui.AnsiOffsetText(40, 0, chat.GetConnChatWindow(conn))
		result += gui.AnsiOffsetText(0, 0, terminal.GetConnTerminal(conn))
		result += gui.AnsiOffsetText(0, 9, tmap.GetMapWindow(conn))
	}

	result += "> "

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

func HandleCommandResponse(conn net.Conn, data parsing.CommandResponse) string {

}
