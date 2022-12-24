package terminal

import (
  "windows"
  "entities"
  "fmt"
  "strings"
  "bufio"
  "os"
)

var TERM_HEIGHT int = 20
var TERM_WIDTH int = 80

type Terminal struct {
  Player entities.Player
  Room entities.Room
  Buffer []string
  Win windows.SimpleWindowInfo
}

func CreateTerminal(player entities.Player, room entities.Room) Terminal {
  window := windows.CreateSimpleWindowInfo(0, 0, TERM_HEIGHT, TERM_WIDTH, "")
  controller := Terminal{player, room, make([]string, TERM_HEIGHT - 3), window}
  return controller
}

func (tc *Terminal) PreDrawFunc() string {
  return windows.BoxedPreDrawFunc(tc)
}

func (tc *Terminal) Clear() {
  windows.DefaultClearFunc(tc)
}

func (tc *Terminal) Draw() {
  windows.DefaultDrawFunc(tc)
}

func (tc Terminal) PostDrawFunc() string {
  return "\033[2A\033[4C"
}

func (tc Terminal) X() int {
  return tc.Win.X
}

func (tc Terminal) Y() int {
  return tc.Win.Y
}

func (tc Terminal) H() int {
  return tc.Win.H
}

func (tc Terminal) W() int {
  return tc.Win.W
}

func (tc Terminal) Body() string {
  var formattedLines []string
  var startingIndex int = 0
  var displayedLineCount = TERM_HEIGHT - 3

  if len(tc.Buffer) < displayedLineCount {
    formattedLines = make([]string, displayedLineCount)

    diff := displayedLineCount - len(tc.Buffer)
    for i := 0; i < diff; i++ {
      formattedLines[i] = ""
    }

    startingIndex = diff
  } else {
    formattedLines = make([]string, len(tc.Buffer))
  }

  for i, line := range tc.Buffer {
    formattedLines[startingIndex + i] =
      fmt.Sprintf("%d: %s", len(tc.Buffer) - i, line)
  }

  if len(formattedLines) > displayedLineCount {
    formattedLines = formattedLines[len(formattedLines) - displayedLineCount:]
  }

  formattedLines = append(formattedLines, "?> ")

  return strings.Join(formattedLines, "\n")
}

func (tc *Terminal) PushLine(line string) {
  tc.Buffer = append(tc.Buffer, line)
}

func (tc *Terminal) Interact() windows.Window {
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  tc.PushLine(scanner.Text())
  return tc
}
