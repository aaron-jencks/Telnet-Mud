package ui

import (
	"fmt"
	"time"
)

// The format for timestamps used for the program
const TimestampFormat = "Mon, Jan _2 2006 @ 03:04:05 PM"

// Prepends the local date and time to the given string
func AddTimestamp(suffix string) string {
	return time.Now().Local().Format(TimestampFormat) + suffix
}

// Prepends the local time to the given string
func AddTime(suffix string) string {
	return time.Now().Local().Format("03:04:05 PM") + suffix
}

// Removes non-visible ascii characters from the given string
func StripIllegalChars(data string) string {
	result := ""
	for _, dc := range data {
		if dc < 32 || dc > 126 {
			continue
		}
		result += fmt.Sprint(dc)
	}
	return result
}

func BoldText(data string) string {
	return CSI("1", "m") + data + CSI("2", "m")
}
