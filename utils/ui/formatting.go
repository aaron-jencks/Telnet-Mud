package ui

import (
	"fmt"
	"strings"
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
	bodyLines := strings.Split(body, "\n")
	for _, bline := range bodyLines {
		for StringLength(bline) > width {
			nearestWord := FindNearestWordBoundaryR(bline, width-1)

			var line string
			if nearestWord <= 0 {
				line = FindFirstNCharacters(bline, width)
				bline = bline[len(line):]
			} else {
				line = bline[:nearestWord]
				bline = bline[nearestWord+1:]
			}

			lines = append(lines, line)
		}
		lines = append(lines, bline)
	}

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

// Determines if the byte at the starting index is a header marker for
// utf-8 and returns how many bytes the character will consume
func IsUnicode(s string, start int) (bool, int) {
	b := s[start]
	if b&0xc0 == 0xc0 {
		// Yes it is
		length := 2
		if b&0x20 > 0 {
			// It has at least 3 bytes
			length++
			if b&0x10 > 0 {
				// It has 4 bytes
				length++
			}
		}

		return true, length
	}
	return false, 1
}

// Returns the length of the string, barring any escape sequences
func StringLength(s string) int {
	count := 0
	tlen := len(s)

	for bi := 0; bi < tlen; bi++ {
		uni, ulen := IsUnicode(s, bi)
		if uni {
			bi += ulen - 1
			count++
			continue
		}

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

func AddBackground(bg int, inner string) string {
	return CSI(fmt.Sprint(bg), "m") + inner + CSI("0", "m")
}
