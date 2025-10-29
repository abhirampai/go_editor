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
var offsetX, offsetY int
var currentX, currentY int
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

func displayText() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		textBufferRow := row + offsetY
		textBufferCol := offsetX
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
			} else if row+offsetY > len(textBuffer) {
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
		{text: getModeModeText(), fg: C.TB_BLACK, bg: C.TB_GREEN, separator: true},
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

func getModeModeText() string {
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
	return fmt.Sprintf("Ln %d, Col %d", currentY+1, currentX+1)
}

func getTabSizeText() string {
	return fmt.Sprintf("Tab Size: %d", editSettings.TabSize)
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

	for {
		COLS = int(C.tb_width())
		ROWS = int(C.tb_height())
		ROWS--
		if COLS < 78 {
			COLS = 78
		}
		C.tb_clear()
		displayText()
		displayStatusBar()
		C.tb_present()
		C.tb_poll_event(&event)
		if event._type == C.TB_EVENT_KEY && event.key == C.TB_KEY_ESC {
			C.tb_shutdown()
			break
		}
	}
}
