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

	"github.com/mattn/go-runewidth"
)

var editSettings, err = LoadSettings()

var sourceFile string
var ROWS, COLS int
var offsetX, offsetY int

var textBuffer = [][]rune{}

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
		C.tb_present()
		C.tb_poll_event(&event)
		if event._type == C.TB_EVENT_KEY && event.key == C.TB_KEY_ESC {
			C.tb_shutdown()
			break
		}
	}
}
