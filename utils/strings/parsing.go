package strings

import "strings"

// Splits a string on the given separator
// Unless the separator appears between two
// Quotes "Like this" or is escaped,
// like\ this
func SplitWithQuotes(s string, sep rune) []string {
	var result []string

	previous := 0
	escaped := false
	quoted := false
	for bi, b := range s {
		if escaped {
			escaped = false

		} else if quoted {
			if b == '"' {
				quoted = false
			} else if b == '\\' {
				escaped = true
			}

		} else {
			if b == '"' {
				quoted = true
			} else if b == '\\' {
				escaped = true

			} else if b == sep {
				if bi != previous {
					result = append(result, s[previous:bi])
				} else {
					result = append(result, "")
				}

				previous = bi + 1
			}
		}
	}

	if previous < len(s) {
		result = append(result, s[previous:])
	} else {
		result = append(result, "")
	}

	return result
}

func StripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func StringContains(arr []string, target string) bool {
	for _, sa := range arr {
		if sa == target {
			return true
		}
	}
	return false
}

func PrepareStringForQuery(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
