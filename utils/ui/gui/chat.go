package gui

import (
	"fmt"
	"mud/utils/ui"
)

func FormatChatEntry(player string, entry string) string {
	return ui.AddTime(fmt.Sprintf(" %s: %s", player, entry))
}
