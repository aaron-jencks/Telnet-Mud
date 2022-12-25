package ui

import "time"

const TimestampFormat = "Mon, Jan _2 2006 @ 03:04:05 PM"

func AddTimestamp(suffix string) string {
	return time.Now().Local().Format(TimestampFormat) + suffix
}
