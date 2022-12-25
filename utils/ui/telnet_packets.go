package ui

import (
	"mud/utils/ui/logger"
	"strings"
)

func ESC(code string) string {
	return "\033" + code
}

func CSI(args ...string) string {
	if len(args) > 1 {
		return "\033[" + strings.Join(args[:len(args)-1], ";") + args[len(args)-1]
	} else if len(args) == 0 {
		logger.Error("Escape code must be supplied with at least 1 argument")
		panic(args)
	} else {
		return "\033[" + args[0]
	}
}

func SaveAndResetCursor(inner string) string {
	return ESC("7") + inner + ESC("8")
}
