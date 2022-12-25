package gui

import (
	"fmt"
	"mud/utils/ui"
	"strings"
)

func Clearscreen() {
	fmt.Print(ui.CSI("2", "J") + ui.CSI(";H"))
}

func CenterAlignText(text string, length int) string {
	tlen := len(text)

	if tlen >= length {
		return text
	}

	diff := length - tlen
	midway := diff / 2

	fmtString1 := strings.Repeat(" ", midway)
	fmtString2 := strings.Repeat(" ", midway)

	if diff%2 == 1 {
		fmtString2 = strings.Repeat(" ", midway+1)
	}

	return fmtString1 + text + fmtString2
}

func LeftAlignText(text string, length int) string {
	tlen := len(text)

	if tlen >= length {
		return text
	}

	diff := length - tlen

	return text + strings.Repeat(" ", diff)
}

func BoxText(text string) string {
	lines := strings.Split(text, "\n")

	var maxLength int = 0
	for _, line := range lines {
		lineLength := len(line)
		if lineLength > maxLength {
			maxLength = lineLength
		}
	}

	bar := strings.Repeat("\u2500", maxLength)
	top := fmt.Sprintf("\u250c%s\u2510", bar)
	bottom := fmt.Sprintf("\u2514%s\u2518", bar)

	result := top + "\n"
	for _, line := range lines {
		result += fmt.Sprintf("\u2502%s\u2502\n",
			LeftAlignText(line, maxLength))
	}
	result += bottom

	return result
}

func TruncateText(text string, length int) string {
	if len(text) > length {
		return text[:length-1] + ">"
	}
	return text
}

func SizedBoxText(text string, h, w int) string {
	lines := strings.Split(text, "\n")
	var formattedLines []string = make([]string, len(lines))
	for i, line := range lines {
		if len(line) > w-2 {
			formattedLines[i] = TruncateText(line, w-2)
		} else if len(line) < w-2 {
			formattedLines[i] = LeftAlignText(line, w-2)
		} else {
			formattedLines[i] = line
		}
	}

	if len(formattedLines) > h-2 {
		formattedLines = formattedLines[:h-2]
	}

	return BoxText(strings.Join(formattedLines, "\n"))
}

func AnsiOffsetText(x, y int, text string) string {
	lines := strings.Split(text, "\n")

	result := ui.CSI(fmt.Sprint(y+1), fmt.Sprint(x+1), "H")

	sep := ui.CSI("B") + ui.CSI(fmt.Sprint(x), "G")

	result += strings.Join(lines, sep) + sep

	return result
}

func CreateMenu(title, prompt string, entries []string, h, w int) string {
	innerText := CenterAlignText(title, w-2) + "\n"
	var formattedStrings []string = make([]string, len(entries))
	for ei, entry := range entries {
		formattedStrings[ei] = fmt.Sprintf("%d: %s", ei+1, entry)
	}

	innerText += strings.Join(formattedStrings, "\n")
	return SizedBoxText(innerText, h-1, w) + "\n" + prompt
}
