package chat

import (
	"mud/services/chat"
	"mud/utils"
	"mud/utils/ui/gui"
	"net"
	"strings"
)

func DisplayChat(entries []string) string {
	var logLines []string

	// Traverse through entries backward
	for ein := range entries {
		ei := len(entries) - ein - 1

		// Create paragraph
		var lines []string
		buffer := entries[ei]
		for len(buffer) > utils.CHAT_W-2 {
			lines = append(lines, buffer[:utils.CHAT_W-1])
			buffer = buffer[utils.CHAT_W-1:]
		}
		lines = append(lines, buffer)

		if len(lines) > utils.CHAT_H-2 {
			// Truncate down to target length
			lines = lines[len(lines)-(utils.CHAT_H-2):]
			logLines = lines
			break
		}

		// Append new text history lines

		logLines = append(logLines, lines...)

		if len(logLines) >= utils.CHAT_H-2 {
			break
		}
	}

	if len(logLines) > utils.CHAT_H-2 {
		logLines = logLines[len(logLines)-(utils.CHAT_H-2):]
	}

	return gui.SizedBoxText(strings.Join(logLines, "\n"), utils.CHAT_H, utils.CHAT_W)
}

func GetConnChatWindow(conn net.Conn) string {
	return DisplayChat(chat.MessageLogMap[conn])
}
