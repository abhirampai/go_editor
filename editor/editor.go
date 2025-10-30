package editor

/*
#cgo CFLAGS: -I${SRCDIR}/termbox2
#cgo LDFLAGS: ${SRCDIR}/termbox2/libtermbox2.a -lm -lc
#include "termbox2.h"
*/
import "C"

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

var editSettings, err = LoadSettings()

var mode int
var sourceFile string
var ROWS, COLS int
var offsetColumn, offsetRow int
var currentColumn, currentRow int
var textBuffer = [][]rune{}
var undoBuffer = [][]rune{}
var copyBuffer = []rune{}
var modified bool

func printCell(col int, row int, fg C.uintattr_t, bg C.uintattr_t, msg string) {
	for _, character := range msg {
		C.tb_set_cell(C.int(col), C.int(row), C.uint32_t(character), fg, bg)
		col += runewidth.RuneWidth(character)
	}
}

func readFile(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		sourceFile = filename
		textBuffer = append(textBuffer, []rune{})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		textBuffer = append(textBuffer, []rune{})

		for i := 0; i < len(line); i++ {
			textBuffer[lineNumber] = append(textBuffer[lineNumber], rune(line[i]))
		}
		lineNumber++
	}
	if lineNumber == 0 {
		textBuffer = append(textBuffer, []rune{})
	}
}

func insertCharacters(keyEvent C.struct_tb_event) {
	insertCharacter := make([]rune, len(textBuffer[currentRow])+1)
	copy(insertCharacter[:currentColumn], textBuffer[currentRow][:currentColumn])
	switch keyEvent.key {
	case C.TB_KEY_SPACE:
		insertCharacter[currentColumn] = rune(' ')
	case C.TB_KEY_TAB:
		insertCharacter[currentColumn] = rune(' ')
	default:
		insertCharacter[currentColumn] = rune(keyEvent.ch)
	}
	copy(insertCharacter[currentColumn+1:], textBuffer[currentRow][currentColumn:])
	textBuffer[currentRow] = insertCharacter
	currentColumn++
	modified = true
}

func deleteCharacter() {
	if currentColumn == 0 && currentRow == 0 {
		return
	}
	if currentColumn > 0 {
		textBuffer[currentRow] = append(textBuffer[currentRow][:currentColumn-1], textBuffer[currentRow][currentColumn:]...)
		currentColumn--
	} else {
		previousRowLength := len(textBuffer[currentRow-1])
		textBuffer[currentRow-1] = append(textBuffer[currentRow-1], textBuffer[currentRow]...)
		textBuffer = append(textBuffer[:currentRow], textBuffer[currentRow+1:]...)
		currentRow--
		currentColumn = previousRowLength
	}
	modified = true
}

func scrollText() {
	if currentRow < offsetRow {
		offsetRow = currentRow
	} else if currentRow >= offsetRow+ROWS {
		offsetRow = currentRow - ROWS + 1
	}

	if currentColumn < offsetColumn {
		offsetColumn = currentColumn
	} else if currentColumn >= offsetColumn+COLS {
		offsetColumn = currentColumn - COLS + 1
	}
}

func displayText() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		textBufferRow := row + offsetRow
		textBufferCol := offsetColumn
		for col = 0; col < COLS; col++ {
			if textBufferRow < len(textBuffer) && textBufferCol < len(textBuffer[textBufferRow]) {
				if textBuffer[textBufferRow][textBufferCol] == '\t' {
					for i := 0; i < editSettings.TabSize; i++ {
						printCell(col, row, C.TB_DEFAULT, C.TB_RED, ".")
						if i < editSettings.TabSize-1 {
							col++
						}
					}
					textBufferCol++
				} else {
					printCell(col, row, C.TB_DEFAULT, C.TB_GREEN, string(textBuffer[textBufferRow][textBufferCol]))
					textBufferCol++
				}
			} else if row+offsetRow > len(textBuffer) {
				printCell(0, row, C.TB_BLUE, C.TB_DEFAULT, "*")
				textBufferCol++
			}
		}
		printCell(col, row, C.TB_DEFAULT, C.TB_DEFAULT, "\n")
	}
}

type statusComponent struct {
	text      string
	fg        C.uintattr_t
	bg        C.uintattr_t
	separator bool
}

func displayStatusBar() {
	const separatorWidth = 2

	copyUndoText, hasCopyUndo := getCopyUndoText()

	leftComponents := []statusComponent{
		{text: getModeStatusText(), fg: C.TB_BLACK, bg: C.TB_GREEN, separator: true},
		{text: getFileStatusText(), fg: C.TB_WHITE, bg: C.TB_BLACK, separator: true},
		{text: copyUndoText, fg: C.TB_WHITE, bg: C.TB_BLACK, separator: hasCopyUndo},
	}

	rightComponents := []statusComponent{
		{text: getCursorStatusText(), fg: C.TB_BLACK, bg: C.TB_CYAN, separator: true},
		{text: getTabSizeText(), fg: C.TB_BLACK, bg: C.TB_CYAN, separator: false},
	}

	leftWidth := 0
	for _, component := range leftComponents {
		leftWidth += len(component.text)
		if component.separator {
			leftWidth += separatorWidth
		}
	}

	rightWidth := 0
	for _, component := range rightComponents {
		rightWidth += len(component.text)
		if component.separator {
			rightWidth += separatorWidth
		}
	}

	currentCol := 0
	for _, component := range leftComponents {
		printCell(currentCol, ROWS, component.fg, component.bg, component.text)
		currentCol += len(component.text)
		if component.separator {
			printCell(currentCol, ROWS, C.TB_WHITE, C.TB_BLACK, "  ")
			currentCol += separatorWidth
		}
	}

	middleSpace := COLS - leftWidth - rightWidth
	if middleSpace > 0 {
		spaces := strings.Repeat(" ", middleSpace)
		printCell(currentCol, ROWS, C.TB_WHITE, C.TB_BLACK, spaces)
		currentCol += middleSpace
	}

	for _, component := range rightComponents {
		printCell(currentCol, ROWS, component.fg, component.bg, component.text)
		currentCol += len(component.text)
		if component.separator {
			printCell(currentCol, ROWS, C.TB_BLACK, C.TB_CYAN, "  ")
			currentCol += separatorWidth
		}
	}
}

func getModeStatusText() string {
	if mode > 0 {
		return "-- INSERT --"
	}
	return "-- VISUAL --"
}

func getFileStatusText() string {
	filenameLength := len(sourceFile)
	if filenameLength > 8 {
		filenameLength = 8
	}
	status := sourceFile[:filenameLength] + " - " + strconv.Itoa(len(textBuffer)) + " lines"
	if modified {
		status += " (modified)"
	} else {
		status += " (saved)"
	}
	return status
}

func getCopyUndoText() (string, bool) {
	var status strings.Builder
	hasContent := false
	if len(copyBuffer) > 0 {
		status.WriteString(" [Copy]")
		hasContent = true
	}
	if len(undoBuffer) > 0 {
		status.WriteString(" [Undo]")
		hasContent = true
	}
	return status.String(), hasContent
}

func getCursorStatusText() string {
	return fmt.Sprintf("Ln %d, Col %d", currentRow+1, currentColumn+1)
}

func getTabSizeText() string {
	return fmt.Sprintf("Tab Size: %d", editSettings.TabSize)
}

func hasTabInRow(row int) bool {
	if row >= len(textBuffer) {
		return false
	}
	for _, ch := range textBuffer[row] {
		if ch == '\t' {
			return true
		}
	}
	return false
}

func textBufferRowLength(row int) int {
	if row >= len(textBuffer) {
		return 0
	}

	if hasTabInRow(row) {
		return len(textBuffer[row]) + (editSettings.TabSize-1)*strings.Count(string(textBuffer[row]), "\t")
	}

	return len(textBuffer[row])
}

func processKeypress(keyEvent C.struct_tb_event) {
	if keyEvent.key == C.TB_KEY_ESC {
		if mode > 0 {
			mode = 0
			return
		}
	} else if keyEvent.ch != 0 {
		if mode > 0 {
			insertCharacters(keyEvent)
		} else {
			switch keyEvent.ch {
			case 'q':
				C.tb_shutdown()
				os.Exit(0)
			case 'i':
				mode = 1
			}
		}
	} else {
		switch keyEvent.key {
		case C.TB_KEY_BACKSPACE, C.TB_KEY_BACKSPACE2:
			if mode > 0 {
				deleteCharacter()
			}
		case C.TB_KEY_TAB:
			if mode > 0 {
				for i := 0; i < editSettings.TabSize; i++ {
					insertCharacters(keyEvent)
				}
			}
		case C.TB_KEY_SPACE:
			if mode > 0 {
				insertCharacters(keyEvent)
			}
		case C.TB_KEY_HOME:
			currentColumn = 0
		case C.TB_KEY_END:
			currentColumn = textBufferRowLength(currentRow)
		case C.TB_KEY_PGUP:
			if currentRow-int(ROWS/4) > 0 {
				currentRow -= int(ROWS / 4)
			}
		case C.TB_KEY_PGDN:
			if currentRow+int(ROWS/4) < len(textBuffer)-1 {
				currentRow += int(ROWS / 4)
			}
		case C.TB_KEY_ARROW_UP:
			if currentRow > 0 {
				currentRow--
			}
		case C.TB_KEY_ARROW_DOWN:
			if currentRow < len(textBuffer)-1 {
				currentRow++
			}
		case C.TB_KEY_ARROW_LEFT:
			if currentColumn > 0 {
				currentColumn--
			} else if currentRow > 0 {
				currentRow--
				currentColumn = textBufferRowLength(currentRow)
			}
		case C.TB_KEY_ARROW_RIGHT:
			if currentColumn < textBufferRowLength(currentRow) {
				currentColumn++
			} else if currentRow < len(textBuffer)-1 {
				currentRow++
				currentColumn = 0
			}
		}
		if currentColumn > textBufferRowLength(currentRow) {
			currentColumn = textBufferRowLength(currentRow)
		}
	}
}

func RunEditor() {
	event := C.struct_tb_event{}

	err := C.tb_init()
	if err != 0 {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		sourceFile = os.Args[1]
		readFile(sourceFile)
	} else {
		sourceFile = "untitled"
		textBuffer = append(textBuffer, []rune{})
	}

	currentRow = 0

	for {
		COLS = int(C.tb_width())
		ROWS = int(C.tb_height())
		ROWS--
		if COLS < 78 {
			COLS = 78
		}
		C.tb_clear()
		scrollText()
		displayText()
		displayStatusBar()
		C.tb_set_cursor(C.int(currentColumn-offsetColumn), C.int(currentRow-offsetRow))
		C.tb_present()
		C.tb_poll_event(&event)
		if event._type == C.TB_EVENT_KEY {
			processKeypress(event)
		}
	}
}
