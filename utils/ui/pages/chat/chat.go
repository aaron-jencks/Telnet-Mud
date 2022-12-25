package chat

import (
	"fmt"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"strings"
)

const (
	CHAT_H = 20
	CHAT_W = 80
)

func FormatChatEntry(player string, entry string) string {
	return ui.AddTime(fmt.Sprintf(" %s: %s", player, entry))
}

func DisplayChat(entries []string) string {
	var logLines []string

	// Traverse through entries backward
	for ein := range entries {
		ei := len(entries) - ein - 1

		// Create paragraph
		var lines []string
		buffer := entries[ei]
		for len(buffer) > CHAT_W-2 {
			lines = append(lines, buffer[:CHAT_W-1])
			buffer = buffer[CHAT_W-1:]
		}
		lines = append(lines, buffer)

		if len(lines) > CHAT_H-2 {
			// Truncate down to target length
			lines = lines[len(lines)-(CHAT_H-2):]
			logLines = lines
			break
		}

		// Append new text history lines

		logLines = append(logLines, lines...)

		if len(logLines) >= CHAT_H-2 {
			break
		}
	}

	if len(logLines) > CHAT_H-2 {
		logLines = logLines[len(logLines)-(CHAT_H-2):]
	}

	return gui.SizedBoxText(strings.Join(logLines, "\n"), CHAT_H, CHAT_W)
}
