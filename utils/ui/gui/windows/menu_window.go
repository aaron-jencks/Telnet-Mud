package windows

import (
	"fmt"
	"mud/utils/ui/gui"
	"mud/utils/ui/logger"
)

type MenuWindow struct {
	Title     string
	Prompt    string
	Entries   []string
	ChoiceMap map[string]func() Window
	Info      SimpleWindowInfo
	HasError  bool
	Error     string
}

func CreateMenuWindow(title, prompt string, entries []string,
	choiceMap map[string]func() Window, x, y, h, w int) MenuWindow {
	return MenuWindow{
		title,
		prompt,
		entries,
		choiceMap,
		CreateSimpleWindowInfo(x, y, h, w, ""),
		false,
		"",
	}
}

func (mw MenuWindow) X() int {
	return mw.Info.X
}

func (mw MenuWindow) Y() int {
	return mw.Info.Y
}

func (mw MenuWindow) H() int {
	return mw.Info.H
}

func (mw MenuWindow) W() int {
	return mw.Info.W
}

func (mw MenuWindow) Body() string {
	result := gui.CreateMenu(mw.Title, mw.Prompt, mw.Entries, mw.H(), mw.W())

	if mw.HasError {
		mw.HasError = false
		result = fmt.Sprintf("\033[91m[ERROR] \033[0m%s\n", mw.Error) + result
	}

	return result
}

func (mw *MenuWindow) Draw() {
	DefaultDrawFunc(mw)
}

func (mw *MenuWindow) Clear() {
	DefaultClearFunc(mw)
}

func (mw *MenuWindow) PreDrawFunc() string {
	return DefaultPreDrawFunc(mw)
}

func (mw *MenuWindow) PostDrawFunc() string {
	return fmt.Sprintf("\033[1A\033[%dC", len(mw.Prompt))
}

func (mw *MenuWindow) Interact() Window {
	var choice int
	fmt.Scanf("%d", &choice)
	if choice <= 0 || choice > len(mw.Entries) {
		mw.HasError = true
		mw.Error = fmt.Sprintf(
			"%d is not a valid choice, please enter a number in (1-%d)",
			choice, len(mw.Entries))
		return mw
	}
	entryText := mw.Entries[choice-1]
	v, ok := mw.ChoiceMap[entryText]
	if ok {
		return v()
	} else {
		logger.Warn("Menu Window Choice Map did not have an entry for %s", entryText)
	}
	return nil
}
