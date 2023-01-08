package terminal

import (
	"fmt"
	"mud/services/terminal"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"net"
	"strings"
)

func DisplayTerminal(ts *terminal.Terminal) string {
	var logLines []string

	// Traverse through ts.Buffer backward
	var ein int = 0
	for ein = range ts.Buffer {
		ei := len(ts.Buffer) - ein - 1

		lines := ui.CreateTextParagraph(
			fmt.Sprintf("%d: %s", ein+1, ts.Buffer[ei]),
			utils.TERMINAL_W-2)

		if len(lines) > utils.TERMINAL_H-2 {
			// Truncate down to target length
			lines = lines[len(lines)-(utils.TERMINAL_H-2):]

			// Reverse lines
			var reversedLines []string = make([]string, len(lines))
			for li, lentry := range lines {
				reversedLines[len(lines)-li-1] = lentry
			}

			logLines = reversedLines
			break
		}

		// Reverse lines
		var reversedLines []string = make([]string, len(lines))
		for li, lentry := range lines {
			reversedLines[len(lines)-li-1] = lentry
		}

		// Append new text history lines
		logLines = append(logLines, reversedLines...)

		if len(logLines) >= utils.TERMINAL_H-2 {
			break
		}
	}

	if len(logLines) > utils.TERMINAL_H-2 {
		logLines = logLines[:utils.TERMINAL_H-2]
	} else if len(logLines) < utils.TERMINAL_H-2 {
		ein = len(ts.Buffer) + 1
		for len(logLines) < utils.TERMINAL_H-2 {
			logLines = append(logLines, fmt.Sprintf("%d: ", ein))
			ein++
		}
	}

	// Reverse Log Lines
	var newLogBuff []string = make([]string, len(logLines))
	for li, lentry := range logLines {
		newLogBuff[len(logLines)-li-1] = lentry
	}

	return gui.SizedBoxText(strings.Join(newLogBuff, "\n"), utils.TERMINAL_H, utils.TERMINAL_W)
}

func GetConnTerminal(conn net.Conn) string {
	return DisplayTerminal(terminal.TerminalMap[conn])
}
