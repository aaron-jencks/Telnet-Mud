package ui

import (
	"fmt"
	"strings"
	"time"
)

const TimestampFormat = "Mon, Jan _2 2006 @ 03:04:05 PM"

func AddTimestamp(suffix string) string {
	return time.Now().Local().Format(TimestampFormat) + suffix
}

func AddTime(suffix string) string {
	return time.Now().Local().Format("03:04:05 PM") + suffix
}

func StripIllegalChars(data string, illegalChars string) string {
	result := ""
	for _, dc := range data {
		if strings.Contains(illegalChars, fmt.Sprint(dc)) {
			continue
		}
		result += fmt.Sprint(dc)
	}
	return result
}
