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
	for StringLength(body) > width {
		nearestWord := FindNearestWordBoundaryR(body, width-1)

		var line string
		if nearestWord <= 0 {
			line = FindFirstNCharacters(body, width)
			body = body[len(line):]
		} else {
			line = body[:nearestWord]
			body = body[nearestWord+1:]
		}

		lines = append(lines, line)
	}
	lines = append(lines, body)

	return lines
}

// Determines if the string contains an Control Sequence
// At the given location, and if so, returns the end index of it
func IsCSI(s string, start int) (bool, int) {
	if s[start] == '\033' {

		// Check for control sequence
		if start < len(s)-1 {
			// There are more bytes left
			current := start
			if s[current+1] == '[' {
				// It's an escape sequence
				current += 2

				// Find the rest of the sequence
				for current < len(s) && (s[current] < 65 || s[current] > 122) {
					current++
				}

				if current >= len(s) {
					return false, -1
				}

				return true, current
			}
		}
	}

	return false, -1
}

// Returns the length of the string, barring any escape sequences
func StringLength(s string) int {
	count := 0
	tlen := len(s)

	for bi := 0; bi < tlen; bi++ {
		csi, clen := IsCSI(s, bi)
		if csi {
			bi = clen
			continue
		}

		count++
	}

	return count
}

// Returns the first n characters of the string,
// barring any escape sequences
func FindFirstNCharacters(s string, n int) string {
	result := ""

	count := 0
	for bi := 0; bi < len(s) && count < n; bi++ {
		csi, clen := IsCSI(s, bi)
		if csi {
			result += s[bi:clen]
			bi = clen
			continue
		}

		result += fmt.Sprintf("%c", s[bi])
		count++
	}

	return result
}
