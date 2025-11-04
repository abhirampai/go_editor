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

type Mode int

const (
	ModeEditor Mode = iota
	ModeHelp
)

var currentMode Mode = ModeEditor

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

func runeIndexToDisplayCol(row int, runeIndex int) int {
	if row >= len(textBuffer) {
		return 0
	}
	col := 0
	if runeIndex > len(textBuffer[row]) {
		runeIndex = len(textBuffer[row])
	}
	for i := 0; i < runeIndex; i++ {
		ch := textBuffer[row][i]
		if ch == '\t' {
			col += editSettings.TabSize
		} else {
			col += runewidth.RuneWidth(ch)
		}
	}
	return col
}

func displayColToRuneIndex(row int, displayCol int) int {
	if row >= len(textBuffer) {
		return 0
	}
	col := 0
	for i := 0; i < len(textBuffer[row]); i++ {
		ch := textBuffer[row][i]
		width := 0
		if ch == '\t' {
			width = editSettings.TabSize
		} else {
			width = runewidth.RuneWidth(ch)
		}
		if col+width > displayCol {
			return i
		}
		col += width
	}
	return len(textBuffer[row])
}

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

func writeFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range textBuffer {
		for _, ch := range line {
			_, err = writer.WriteRune(ch)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		writer.WriteRune('\n')
	}
	writer.Flush()
	modified = false
}

func insertCharacters(keyEvent C.struct_tb_event) {
	insertCharacter := make([]rune, len(textBuffer[currentRow])+1)
	copy(insertCharacter[:currentColumn], textBuffer[currentRow][:currentColumn])
	switch keyEvent.key {
	case C.TB_KEY_SPACE:
		insertCharacter[currentColumn] = rune(' ')
	case C.TB_KEY_TAB:
		insertCharacter[currentColumn] = rune('\t')
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

func insertNewLine() {
	newLine := make([]rune, len(textBuffer[currentRow])-currentColumn)
	copy(newLine, textBuffer[currentRow][currentColumn:])
	textBuffer[currentRow] = textBuffer[currentRow][:currentColumn]
	textBuffer = append(textBuffer[:currentRow+1], textBuffer[currentRow:]...)
	textBuffer[currentRow+1] = newLine
	currentRow++
	currentColumn = 0
	modified = true
}

func copyLine() {
	if currentRow < len(textBuffer) {
		copyBuffer = make([]rune, len(textBuffer[currentRow]))
		copy(copyBuffer, textBuffer[currentRow])
	}
}

func pasteLine() {
	if len(copyBuffer) == 0 {
		currentRow++
		currentColumn = 0
	}
	insertedLine := make([]rune, len(copyBuffer))
	copy(insertedLine, copyBuffer)
	textBuffer = append(textBuffer[:currentRow+1], textBuffer[currentRow:]...)
	textBuffer[currentRow] = insertedLine
	modified = true
}

func deleteLine() {
	copyLine()
	if currentRow < len(textBuffer) {
		textBuffer = append(textBuffer[:currentRow], textBuffer[currentRow+1:]...)
		if currentRow >= len(textBuffer) && currentRow > 0 {
			currentRow--
		}
		currentColumn = 0
		modified = true
	}
}

func pushBuffer() {
	copyUndoBuffer := make([][]rune, len(textBuffer))
	copy(copyUndoBuffer, textBuffer)
	undoBuffer = copyUndoBuffer
}

func pullBuffer() {
	if len(undoBuffer) == 0 {
		return
	}

	textBuffer = undoBuffer
	undoBuffer = [][]rune{}
}

func scrollText() {
	if currentRow < offsetRow {
		offsetRow = currentRow
	} else if currentRow >= offsetRow+ROWS {
		offsetRow = currentRow - ROWS + 1
	}

	visCol := 0
	if currentRow < len(textBuffer) {
		visCol = runeIndexToDisplayCol(currentRow, currentColumn)
	}
	if visCol < offsetColumn {
		offsetColumn = visCol
	} else if visCol >= offsetColumn+COLS {
		offsetColumn = visCol - COLS + 1
	}
}

func displayText() {
	for scrRow := 0; scrRow < ROWS; scrRow++ {
		textRow := scrRow + offsetRow
		if textRow >= len(textBuffer) {
			printCell(0, scrRow, C.TB_BLUE, C.TB_DEFAULT, "*")
			continue
		}

		startRune := displayColToRuneIndex(textRow, offsetColumn)
		visCol := runeIndexToDisplayCol(textRow, startRune)

		for runeIdx := startRune; runeIdx < len(textBuffer[textRow]) && visCol-offsetColumn < COLS; runeIdx++ {
			ch := textBuffer[textRow][runeIdx]
			if ch == '\t' {
				for i := 0; i < editSettings.TabSize && visCol-offsetColumn < COLS; i++ {
					printCell(visCol-offsetColumn, scrRow, C.TB_DEFAULT, C.TB_RED, ".")
					visCol++
				}
			} else {
				printCell(visCol-offsetColumn, scrRow, C.TB_DEFAULT, C.TB_GREEN, string(ch))
				visCol += runewidth.RuneWidth(ch)
			}
		}
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
	visCol := 0
	if currentRow < len(textBuffer) {
		visCol = runeIndexToDisplayCol(currentRow, currentColumn)
	}
	return fmt.Sprintf("Ln %d, Col %d", currentRow+1, visCol+1)
}

func getTabSizeText() string {
	return fmt.Sprintf("Tab Size: %d", editSettings.TabSize)
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
			case 'w':
				writeFile(sourceFile)
			case 'h':
				currentMode = ModeHelp
			case 'c':
				copyLine()
			case 'p':
				pasteLine()
			case 'd':
				deleteLine()
			case 's':
				pushBuffer()
			case 'l':
				pullBuffer()
			}
		}
	} else {
		switch keyEvent.key {
		case C.TB_KEY_ENTER:
			if mode > 0 {
				insertNewLine()
			}
		case C.TB_KEY_BACKSPACE, C.TB_KEY_BACKSPACE2:
			if mode > 0 {
				deleteCharacter()
			}
		case C.TB_KEY_TAB:
			if mode > 0 {
				insertCharacters(keyEvent)
			}
		case C.TB_KEY_SPACE:
			if mode > 0 {
				insertCharacters(keyEvent)
			}
		case C.TB_KEY_HOME:
			currentColumn = 0
		case C.TB_KEY_END:
			if currentRow < len(textBuffer) {
				currentColumn = len(textBuffer[currentRow])
			} else {
				currentColumn = 0
			}
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
				currentColumn = len(textBuffer[currentRow])
			}
		case C.TB_KEY_ARROW_RIGHT:
			if currentRow < len(textBuffer) && currentColumn < len(textBuffer[currentRow]) {
				currentColumn++
			} else if currentRow < len(textBuffer)-1 {
				currentRow++
				currentColumn = 0
			}
		}
		if currentRow < len(textBuffer) && currentColumn > len(textBuffer[currentRow]) {
			currentColumn = len(textBuffer[currentRow])
		}
	}
}

func drawPopupFrame(x, y, w, h int, title string) {
	for i := 0; i < w; i++ {
		C.tb_set_cell(C.int(x+i), C.int(y), '─', C.TB_WHITE, C.TB_BLACK)
		C.tb_set_cell(C.int(x+i), C.int(y+h-1), '─', C.TB_WHITE, C.TB_BLACK)
	}
	for j := 0; j < h; j++ {
		C.tb_set_cell(C.int(x), C.int(y+j), '│', C.TB_WHITE, C.TB_BLACK)
		C.tb_set_cell(C.int(x+w-1), C.int(y+j), '│', C.TB_WHITE, C.TB_BLACK)
	}

	C.tb_set_cell(C.int(x), C.int(y), '┌', C.TB_WHITE, C.TB_BLACK)
	C.tb_set_cell(C.int(x+w-1), C.int(y), '┐', C.TB_WHITE, C.TB_BLACK)
	C.tb_set_cell(C.int(x), C.int(y+h-1), '└', C.TB_WHITE, C.TB_BLACK)
	C.tb_set_cell(C.int(x+w-1), C.int(y+h-1), '┘', C.TB_WHITE, C.TB_BLACK)

	printCell(x+2, y, C.TB_YELLOW, C.TB_BLACK, title)

	for i := 1; i < w-1; i++ {
		for j := 1; j < h-1; j++ {
			C.tb_set_cell(C.int(x+i), C.int(y+j), ' ', C.TB_DEFAULT, C.TB_BLACK)
		}
	}
}

func showHelp() {
	w := int(C.tb_width())
	h := int(C.tb_height())

	helpText := FormatKeyBindingsHelp()

	maxWidth := 0
	for _, line := range helpText {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	pw := maxWidth + 4      // Add padding for borders
	ph := len(helpText) + 4 // Add space for borders and footer

	x := (w - pw) / 2
	y := (h - ph) / 2

	drawPopupFrame(x, y, pw, ph, "Help")

	for i, line := range helpText {
		printCell(x+2, y+1+i, C.TB_WHITE, C.TB_BLACK, line)
	}

	footerText := "[Enter/Esc] Close"
	footerX := x + (pw-len(footerText))/2
	printCell(footerX, y+ph-2, C.TB_BLUE, C.TB_BLACK, footerText)

	C.tb_present()
}

func processPopover(event C.struct_tb_event) {
	if event.key == C.TB_KEY_ENTER || event.key == C.TB_KEY_ESC {
		currentMode = ModeEditor
		return
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

		switch currentMode {
		case ModeEditor:
			scrollText()
			displayText()
			displayStatusBar()
			visCol := 0
			if currentRow < len(textBuffer) {
				visCol = runeIndexToDisplayCol(currentRow, currentColumn)
			}
			C.tb_set_cursor(C.int(visCol-offsetColumn), C.int(currentRow-offsetRow))
		case ModeHelp:
			displayText()
			displayStatusBar()
			showHelp()
			C.tb_set_cursor(-1, -1)
		}

		C.tb_present()
		C.tb_poll_event(&event)
		if event._type == C.TB_EVENT_KEY {
			switch currentMode {
			case ModeEditor:
				processKeypress(event)
			case ModeHelp:
				processPopover(event)
			}
		}
	}
}
