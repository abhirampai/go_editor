package main

/*
#cgo CFLAGS: -I${SRCDIR}/termbox2
#cgo LDFLAGS: ${SRCDIR}/termbox2/libtermbox2.a -lm -lc
#include "termbox2.h"
*/
import "C"

import (
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
)

var ROWS, COLS int
var offsetX, offsetY int

var textBuffer = [][]rune{
	{'H', 'e', 'l', 'l', 'o'},
	{'W', 'o', 'r', 'l', 'd'},
}

func printCell(col int, row int, fg C.uintattr_t, bg C.uintattr_t, msg string) {
	for _, character := range msg {
		C.tb_set_cell(C.int(col), C.int(row), C.uint32_t(character), fg, bg)
		col += runewidth.RuneWidth(character)
	}
}

func displayText() {
	var row, col int
	for row = 0; row < ROWS; row++ {
		textBufferRow := row + offsetY
		for col = 0; col < COLS; col++ {
			textBufferCol := col + offsetX
			if textBufferRow < len(textBuffer) && textBufferCol < len(textBuffer[textBufferRow]) {
				if textBuffer[textBufferRow][textBufferCol] == '\t' {
					printCell(col, row, C.TB_DEFAULT, C.TB_DEFAULT, " ")
				} else {
					printCell(col, row, C.TB_DEFAULT, C.TB_GREEN, string(textBuffer[textBufferRow][textBufferCol]))
				}
			} else if row+offsetY > len(textBuffer) {
				printCell(0, row, C.TB_BLUE, C.TB_DEFAULT, "*")
			}
		}
		printCell(col, row, C.TB_DEFAULT, C.TB_DEFAULT, "\n")
	}
}

func runEditor() {
	event := C.struct_tb_event{}

	err := C.tb_init()
	if err != 0 {
		fmt.Println(err)
		os.Exit(1)
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

func main() {
	runEditor()
}
