package chat

import (
	"mud/services/chat"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"net"
	"strings"
)

func DisplayChat(entries []string) string {
	var logLines []string

	// Traverse through entries backward
	for ein := range entries {
		ei := len(entries) - ein - 1

		lines := ui.CreateTextParagraph(entries[ei], utils.CHAT_W-2)

		if len(lines) > utils.CHAT_H-2 {
			// Truncate down to target length
			lines = lines[len(lines)-(utils.CHAT_H-2):]
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

		if len(logLines) >= utils.CHAT_H-2 {
			break
		}
	}

	if len(logLines) > utils.CHAT_H-2 {
		logLines = logLines[len(logLines)-(utils.CHAT_H-2):]
	}

	// Reverse Log Lines
	var newLogBuff []string = make([]string, len(logLines))
	for li, lentry := range logLines {
		newLogBuff[len(logLines)-li-1] = lentry
	}

	return gui.SizedBoxText(strings.Join(newLogBuff, "\n"), utils.CHAT_H, utils.CHAT_W) + "\n> "
}

func GetConnChatWindow(conn net.Conn) string {
	return DisplayChat(chat.MessageLogMap[conn])
}
