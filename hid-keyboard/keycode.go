package main

type Keycode uint16

const (
	NumRows   = 15
	NumCols   = 7
	NumLayers = 2

	KC_NO         Keycode = 0x0000
	KC_BOOTLOADER Keycode = 0xF000 // [13,5] QK_BOOT → machine.EnterBootloader
)

func idx(row, col int) int { return row*NumCols + col }
