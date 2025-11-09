package editor

/*
#include "termbox2.h"
*/
import "C"

func showFileBrowser() {
	w := int(C.tb_width())
	h := int(C.tb_height())

	pw := w - 20
	ph := h - 6

	x := (w - pw) / 2
	y := (h - ph) / 2

	drawPopupFrame(x, y, pw, ph, "File Browser")

	displayPath := fileBrowser.CurrentPath
	if len(displayPath) > pw-6 {
		displayPath = "..." + displayPath[len(displayPath)-(pw-9):]
	}
	printCell(x+2, y+1, C.TB_WHITE, C.TB_BLACK, displayPath)

	visibleEntries := ph - 4
	startIdx := fileBrowser.Scroll
	endIdx := min(startIdx+visibleEntries, len(fileBrowser.Entries))

	for i := startIdx; i < endIdx; i++ {
		entry := fileBrowser.Entries[i]
		displayName := entry.Name
		if entry.IsDir {
			displayName = "[" + displayName + "]"
		}

		if len(displayName) > pw-6 {
			displayName = displayName[:pw-9] + "..."
		}

		var fg, bg C.uintattr_t = C.TB_WHITE, C.TB_BLACK

		if i == fileBrowser.Cursor {
			fg = C.TB_BLACK
			bg = C.TB_WHITE
		}

		printCell(x+2, y+2+(i-startIdx), fg, bg, displayName)
	}

	footerText := "[↑/↓] Navigate  [Enter] Select  [Esc] Close"
	footerX := x + (pw-len(footerText))/2
	printCell(footerX, y+ph-2, C.TB_BLUE, C.TB_BLACK, footerText)

	C.tb_present()
}

type themeEntry struct {
	Key  string
	Name string
}

var themeSelector = struct {
	Entries []themeEntry
	Cursor  int
}{}

var previousThemeKey string

func showThemeSelector() {
	w := int(C.tb_width())
	h := int(C.tb_height())

	pw := 60
	ph := 8 + len(themeSelector.Entries)
	if ph > h-4 {
		ph = h - 4
	}

	x := (w - pw) / 2
	y := (h - ph) / 2
	if previousThemeKey == "" {
		if editSettings != nil && editSettings.Theme != "" {
			previousThemeKey = editSettings.Theme
		} else {
			previousThemeKey = GetCurrentThemeKey()
		}

		for i, e := range themeSelector.Entries {
			if e.Key == previousThemeKey {
				themeSelector.Cursor = i
				SetTheme(e.Key)
				break
			}
		}

		if len(themeSelector.Entries) > 0 {
			SetTheme(themeSelector.Entries[themeSelector.Cursor].Key)
		}
	}

	drawPopupFrame(x, y, pw, ph, "Select Theme")

	for i, entry := range themeSelector.Entries {
		displayName := entry.Name
		if len(displayName) > pw-6 {
			displayName = displayName[:pw-9] + "..."
		}
		var fg, bg C.uintattr_t = C.TB_WHITE, C.TB_BLACK
		if i == themeSelector.Cursor {
			fg, bg = C.TB_BLACK, C.TB_WHITE
		}
		printCell(x+2, y+2+i, fg, bg, displayName)
	}

	footerText := "[↑/↓] Navigate  [Enter] Select  [Esc] Close"
	footerX := x + (pw-len(footerText))/2
	printCell(footerX, y+ph-2, C.TB_BLUE, C.TB_BLACK, footerText)

	C.tb_present()
}

func processThemeSelectorEvent(event C.struct_tb_event) {
	switch event.key {
	case C.TB_KEY_ESC:
		if previousThemeKey != "" {
			SetTheme(previousThemeKey)
		}
		previousThemeKey = ""
		currentMode = ModeEditor
	case C.TB_KEY_ARROW_UP:
		if themeSelector.Cursor > 0 {
			themeSelector.Cursor--
			SetTheme(themeSelector.Entries[themeSelector.Cursor].Key)
		}
	case C.TB_KEY_ARROW_DOWN:
		if themeSelector.Cursor < len(themeSelector.Entries)-1 {
			themeSelector.Cursor++
			SetTheme(themeSelector.Entries[themeSelector.Cursor].Key)
		}
	case C.TB_KEY_ENTER:
		entry := themeSelector.Entries[themeSelector.Cursor]
		SetThemeAndSave(entry.Key)
		previousThemeKey = ""
		currentMode = ModeEditor
	}
}
func processFileBrowserEvent(event C.struct_tb_event) {
	switch event.key {
	case C.TB_KEY_ESC:
		currentMode = ModeEditor
	case C.TB_KEY_ARROW_UP:
		fileBrowser.MoveUp()
	case C.TB_KEY_ARROW_DOWN:
		fileBrowser.MoveDown()
	case C.TB_KEY_ENTER:
		if selectedPath, isDir, err := fileBrowser.Enter(); err == nil {
			if !isDir {
				currentMode = ModeEditor
				saveCurrentFileIfModified()
				sourceFile = selectedPath
				readFile(sourceFile)
				currentRow = 0
				currentColumn = 0
				offsetRow = 0
				offsetColumn = 0
			}
		}
	}
}

func saveCurrentFileIfModified() {
	if modified {
		writeFile(sourceFile)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
