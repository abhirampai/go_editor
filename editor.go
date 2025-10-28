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

func print_message(col int, row int, fg C.uintattr_t, bg C.uintattr_t, msg string) {
	for _, c := range msg {
		C.tb_set_cell(C.int(col), C.int(row), C.uint32_t(c), fg, bg)
		col += runewidth.RuneWidth(c)
	}
}

func run_editor() {
	event := C.struct_tb_event{}

	err := C.tb_init()
	if err != 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		print_message(25, 11, C.TB_DEFAULT, C.TB_DEFAULT, "Go-Editor - A bare bone text editor")
		C.tb_present()
		C.tb_poll_event(&event)
		if event._type == C.TB_EVENT_KEY && event.key == C.TB_KEY_ESC {
			C.tb_shutdown()
			break
		}
	}
}

func main() {
	run_editor()
}
