package ui

import (
	"mud/controllers/telnet/tx"
	"mud/entities"
	"mud/parsing_services/parsing"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"mud/utils/ui/pages/chat"
	"mud/utils/ui/pages/terminal"
	"mud/utils/ui/pages/tmap"
	"net"
	"strings"
)

func GetDisplayForConn(conn net.Conn, p entities.Player,
	loggedIn, saveCursor, clearScreen, wchat, info, dtmap bool) string {

	var result string

	if clearScreen {
		result += gui.Clearscreen()
	} else {
		result += gui.ResetCursorPosition()
	}

	if !loggedIn {
		result += MOTD()
		result += chat.GetConnChatWindowModHeight(conn, utils.CHAT_H-2)
	} else {
		if wchat {
			result += gui.AnsiOffsetText(40, 0, chat.GetConnChatWindow(conn))
		}
		if info {
			result += gui.AnsiOffsetText(0, 0, terminal.GetConnTerminal(conn))
		}
		if dtmap {
			result += gui.AnsiOffsetText(0, 9, tmap.GetMapWindow(p))
		}
	}

	result += gui.AnsiOffsetText(0, 19, "> ")
	if !saveCursor {
		result += gui.EraseToEndOfLine()
	}

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

func HandleCommandResponse(data parsing.CommandResponse) {
	if data.Global {
		for _, client := range tx.Clients {
			tx.SendTarget([]byte(GetDisplayForConn(client, data.Player, data.LoggedIn,
				true, false,
				data.Chat, data.Info, data.Map)), client)
		}
	} else {
		if len(data.Specific) > 0 {
			for _, user := range data.Specific {
				tx.SendTarget([]byte(GetDisplayForConn(user, data.Player, data.LoggedIn,
					true, false,
					data.Chat, data.Info, data.Map)), user)
			}
		}
	}

	tx.SendTarget([]byte(GetDisplayForConn(data.Conn, data.Player, data.LoggedIn,
		data.SaveCursor, data.Clear,
		data.Chat, data.Info, data.Map)), data.Conn)
}
