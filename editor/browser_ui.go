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
