package chat

import (
	"mud/services/chat"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"net"
	"strings"
)

func DisplayChat(entries []string, height int) string {
	var logLines []string

	// Traverse through entries backward
	for ein := range entries {
		ei := len(entries) - ein - 1

		lines := ui.CreateTextParagraph(entries[ei], utils.CHAT_W-2)

		if len(lines) > height-2 {
			// Truncate down to target length
			lines = lines[len(lines)-(height-2):]
			logLines = lines
			break
		}

		// Reverse lines
		var reversedLines []string = make([]string, len(lines))
		for li, lentry := range lines {
			reversedLines[len(lines)-li-1] = lentry
		}

		// Append new text history lines
		logLines = append(logLines, reversedLines...)

		if len(logLines) >= height-2 {
			break
		}
	}

	if len(logLines) > height-2 {
		logLines = logLines[len(logLines)-(height-2):]
	}

	// Reverse Log Lines
	var newLogBuff []string = make([]string, len(logLines))
	for li, lentry := range logLines {
		newLogBuff[len(logLines)-li-1] = lentry
	}

	return gui.SizedBoxText(strings.Join(newLogBuff, "\n"), height, utils.CHAT_W)
}

func GetConnChatWindow(conn net.Conn) string {
	return DisplayChat(chat.MessageLogMap[conn], utils.CHAT_H)
}

func GetConnChatWindowModHeight(conn net.Conn, height int) string {
	return DisplayChat(chat.MessageLogMap[conn], height)
}
