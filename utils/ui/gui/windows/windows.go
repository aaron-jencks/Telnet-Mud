package windows

import (
	"fmt"
	"mud/utils/ui/gui"
	"strings"
)

type SimpleWindowInfo struct {
	X    int
	Y    int
	H    int
	W    int
	Body string
}

func CreateSimpleWindowInfo(x, y, h, w int, body string) SimpleWindowInfo {
	return SimpleWindowInfo{x, y, h, w, body}
}

type Window interface {
	PreDrawFunc() string
	PostDrawFunc() string
	Interact() Window
	Clear()
	Draw()
	X() int
	Y() int
	H() int
	W() int
	Body() string
}

func DefaultPreDrawFunc(w Window) string {
	return w.Body()
}

func DefaultPostDrawFunc(w Window) string {
	return ""
}

func BoxedPreDrawFunc(w Window) string {
	return gui.SizedBoxText(w.Body(), w.H(), w.W())
}

func DefaultDrawFunc(w Window) {
	fmt.Print(gui.AnsiOffsetText(w.X(), w.Y(), w.PreDrawFunc()))
	fmt.Print(w.PostDrawFunc())
}

func DefaultClearFunc(w Window) {
	var lines []string = make([]string, w.H())
	for i := range lines {
		lines[i] = strings.Repeat(" ", w.W())
	}
	fmt.Print(strings.Join(lines, "\n"))
}

type WindowController struct {
	Windows []Window
}

func (wc WindowController) Display() {
	w := wc.Windows[len(wc.Windows)-1]
	w.Clear()
	w.Draw()
}

func (wc *WindowController) Interact() {
	wc.Display()
	w := wc.Windows[len(wc.Windows)-1]
	nw := w.Interact()
	if nw == nil {
		wc.PopWindow()
	} else if nw != w {
		wc.PushWindow(nw)
	}
}

func (wc *WindowController) PushWindow(w Window) {
	wc.Windows = append(wc.Windows, w)
}

func (wc *WindowController) PopWindow() int {
	if len(wc.Windows) > 0 {
		wc.Windows = wc.Windows[:len(wc.Windows)-1]
	}

	return len(wc.Windows)
}

func CreateWindowController() WindowController {
	return WindowController{
		[]Window{},
	}
}
