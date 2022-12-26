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
		result += fmt.Sprintf("%c", dc)
	}
	return result
}

// Enboldens selected text
func BoldText(data string) string {
	return CSI("1", "m") + data + CSI("2", "m")
}

// Finds the nearest word boundary,
// starting from start and working backwards towards
// the beginning of the string
func FindNearestWordBoundaryR(body string, start int) int {
	for i := start; i >= 0; i-- {
		if body[i] == ' ' {
			return i
		}
	}
	return -1
}

// Takes a long line of text and breaks it into several lines,
// using word boundaries if possible
func CreateTextParagraph(body string, width int) []string {
	// Create paragraph
	var lines []string
	for len(body) > width {
		nearestWord := FindNearestWordBoundaryR(body, width-1)

		var line string
		if nearestWord <= 0 {
			line = body[:width]
			body = body[width:]
		} else {
			line = body[:nearestWord]
			body = body[nearestWord+1:]
		}

		lines = append(lines, line)
	}
	lines = append(lines, body)

	return lines
}
