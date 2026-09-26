package main

import (
	"context"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
)

// Same pins and row*7+col order as hid-keyboard-g1.
var cols = []machine.Pin{
	machine.D18, machine.D14, machine.D15,
	machine.D20, machine.D22, machine.D19, machine.D6,
}

var rows = []machine.Pin{
	machine.D8, machine.D9, machine.D10, machine.D11, machine.D7,
	machine.D16, machine.D5, machine.D3, machine.D4, machine.D1,
	machine.D0, machine.D2, machine.D17, machine.D23, machine.D21,
}

// QMK kinesis/keymaps/default_pretty (my-customize), US HID usages.
// Dynamic macros are KC_NO.
var layer0 = [15 * 7]keyboard.Keycode{
	KeyEqual, KeyTab, KeyLeftCtrl, KeyLeftShift, kcNO, kcNO, kcNO,
	Key1, KeyQ, KeyA, KeyZ, KeyGrave, kcNO, kcNO,
	Key2, KeyW, KeyS, KeyX, KeyInsert, KeyLeftAlt, kcNO,
	Key3, KeyE, KeyD, KeyC, KeyLeft, KeyEscape, KeySpace,
	Key4, KeyR, KeyF, KeyV, kcNO, kcNO, KeyWindows,
	Key5, KeyT, KeyG, KeyB, KeyRight, kcLCAG, KeyLeftCtrl,
	Key6, KeyY, KeyH, KeyN, KeyDown, KeyEnter, kcMEH,
	Key7, KeyU, KeyJ, KeyM, kcNO, KeyRightGUI, kcNO,
	Key8, KeyI, KeyK, KeyComma, KeyUp, KeyPageUp, KeyBackspace,
	Key9, KeyO, KeyL, KeyPeriod, KeyLeftBracket, kcNO, KeyPageDown,
	Key0, KeyP, KeySemicolon, KeySlash, KeyRightBracket, kcNO, kcNO,
	KeyMinus, KeyBackslash, KeyQuote, KeyRightShift, kcNO, kcNO, kcNO,
	KeyEscape, KeyF3, KeyF6, KeyF9, KeyF12, kcNO, kcNO,
	KeyF1, KeyF4, KeyF7, KeyF10, kcNO, kcNO, kcNO,
	KeyF2, KeyF5, KeyF8, KeyF11, kcNO, kcNO, kcNO,
}

func main() {
	usb.Product = "tinygo-kint-kb500"

	d := keyboard.New()
	// Index is row*7+col, same order as hid-keyboard-g1.
	d.AddMatrixKeyboard(cols, rows, [][]keyboard.Keycode{layer0[:]})

	loadKeyboardDef()
	if err := d.Loop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
