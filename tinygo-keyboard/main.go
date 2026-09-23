package main

import (
	"context"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes"
)

const (
	kcNO keyboard.Keycode = 0
	// QMK KC_MEH and LCA(KC_LGUI). The low byte of KC_MEH is KC_NO.
	kcMEH  keyboard.Keycode = keycodes.TypeXCtl | keycodes.TypeXSft | keycodes.TypeXAlt
	kcLCAG keyboard.Keycode = keycodes.TypeXCtl | keycodes.TypeXAlt | 0xE3
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
	keycodes.KeyEqual, keycodes.KeyTab, keycodes.KeyLeftCtrl, keycodes.KeyLeftShift, kcNO, kcNO, kcNO,
	keycodes.Key1, keycodes.KeyQ, keycodes.KeyA, keycodes.KeyZ, keycodes.KeyGrave, kcNO, kcNO,
	keycodes.Key2, keycodes.KeyW, keycodes.KeyS, keycodes.KeyX, keycodes.KeyInsert, keycodes.KeyLeftAlt, kcNO,
	keycodes.Key3, keycodes.KeyE, keycodes.KeyD, keycodes.KeyC, keycodes.KeyLeft, keycodes.KeyEscape, keycodes.KeySpace,
	keycodes.Key4, keycodes.KeyR, keycodes.KeyF, keycodes.KeyV, kcNO, kcNO, keycodes.KeyWindows,
	keycodes.Key5, keycodes.KeyT, keycodes.KeyG, keycodes.KeyB, keycodes.KeyRight, kcLCAG, keycodes.KeyLeftCtrl,
	keycodes.Key6, keycodes.KeyY, keycodes.KeyH, keycodes.KeyN, keycodes.KeyDown, keycodes.KeyEnter, kcMEH,
	keycodes.Key7, keycodes.KeyU, keycodes.KeyJ, keycodes.KeyM, kcNO, keycodes.KeyRightGUI, kcNO,
	keycodes.Key8, keycodes.KeyI, keycodes.KeyK, keycodes.KeyComma, keycodes.KeyUp, keycodes.KeyPageUp, keycodes.KeyBackspace,
	keycodes.Key9, keycodes.KeyO, keycodes.KeyL, keycodes.KeyPeriod, keycodes.KeyLeftBracket, kcNO, keycodes.KeyPageDown,
	keycodes.Key0, keycodes.KeyP, keycodes.KeySemicolon, keycodes.KeySlash, keycodes.KeyRightBracket, kcNO, kcNO,
	keycodes.KeyMinus, keycodes.KeyBackslash, keycodes.KeyQuote, keycodes.KeyRightShift, kcNO, kcNO, kcNO,
	keycodes.KeyEscape, keycodes.KeyF3, keycodes.KeyF6, keycodes.KeyF9, keycodes.KeyF12, kcNO, kcNO,
	keycodes.KeyF1, keycodes.KeyF4, keycodes.KeyF7, keycodes.KeyF10, kcNO, kcNO, kcNO,
	keycodes.KeyF2, keycodes.KeyF5, keycodes.KeyF8, keycodes.KeyF11, kcNO, kcNO, kcNO,
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
