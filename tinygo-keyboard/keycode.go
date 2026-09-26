package main

import "github.com/sago35/tinygo-keyboard/keycodes"

// HID Keyboard/Keypad page (0x07). Values match QMK KC_* on a US keyboard.
const (
	KeyA = keycodes.TypeNormal | 0x04
	KeyB = keycodes.TypeNormal | 0x05
	KeyC = keycodes.TypeNormal | 0x06
	KeyD = keycodes.TypeNormal | 0x07
	KeyE = keycodes.TypeNormal | 0x08
	KeyF = keycodes.TypeNormal | 0x09
	KeyG = keycodes.TypeNormal | 0x0A
	KeyH = keycodes.TypeNormal | 0x0B
	KeyI = keycodes.TypeNormal | 0x0C
	KeyJ = keycodes.TypeNormal | 0x0D
	KeyK = keycodes.TypeNormal | 0x0E
	KeyL = keycodes.TypeNormal | 0x0F
	KeyM = keycodes.TypeNormal | 0x10
	KeyN = keycodes.TypeNormal | 0x11
	KeyO = keycodes.TypeNormal | 0x12
	KeyP = keycodes.TypeNormal | 0x13
	KeyQ = keycodes.TypeNormal | 0x14
	KeyR = keycodes.TypeNormal | 0x15
	KeyS = keycodes.TypeNormal | 0x16
	KeyT = keycodes.TypeNormal | 0x17
	KeyU = keycodes.TypeNormal | 0x18
	KeyV = keycodes.TypeNormal | 0x19
	KeyW = keycodes.TypeNormal | 0x1A
	KeyX = keycodes.TypeNormal | 0x1B
	KeyY = keycodes.TypeNormal | 0x1C
	KeyZ = keycodes.TypeNormal | 0x1D

	Key1 = keycodes.TypeNormal | 0x1E
	Key2 = keycodes.TypeNormal | 0x1F
	Key3 = keycodes.TypeNormal | 0x20
	Key4 = keycodes.TypeNormal | 0x21
	Key5 = keycodes.TypeNormal | 0x22
	Key6 = keycodes.TypeNormal | 0x23
	Key7 = keycodes.TypeNormal | 0x24
	Key8 = keycodes.TypeNormal | 0x25
	Key9 = keycodes.TypeNormal | 0x26
	Key0 = keycodes.TypeNormal | 0x27

	KeyEnter     = keycodes.TypeNormal | 0x28
	KeyEscape    = keycodes.TypeNormal | 0x29
	KeyBackspace = keycodes.TypeNormal | 0x2A
	KeyTab       = keycodes.TypeNormal | 0x2B
	KeySpace     = keycodes.TypeNormal | 0x2C
	KeyMinus     = keycodes.TypeNormal | 0x2D
	KeyEqual     = keycodes.TypeNormal | 0x2E

	KeyLeftBracket  = keycodes.TypeNormal | 0x2F
	KeyRightBracket = keycodes.TypeNormal | 0x30
	KeyBackslash    = keycodes.TypeNormal | 0x31
	KeySemicolon    = keycodes.TypeNormal | 0x33
	KeyQuote        = keycodes.TypeNormal | 0x34
	KeyGrave        = keycodes.TypeNormal | 0x35
	KeyComma        = keycodes.TypeNormal | 0x36
	KeyPeriod       = keycodes.TypeNormal | 0x37
	KeySlash        = keycodes.TypeNormal | 0x38

	KeyF1  = keycodes.TypeNormal | 0x3A
	KeyF2  = keycodes.TypeNormal | 0x3B
	KeyF3  = keycodes.TypeNormal | 0x3C
	KeyF4  = keycodes.TypeNormal | 0x3D
	KeyF5  = keycodes.TypeNormal | 0x3E
	KeyF6  = keycodes.TypeNormal | 0x3F
	KeyF7  = keycodes.TypeNormal | 0x40
	KeyF8  = keycodes.TypeNormal | 0x41
	KeyF9  = keycodes.TypeNormal | 0x42
	KeyF10 = keycodes.TypeNormal | 0x43
	KeyF11 = keycodes.TypeNormal | 0x44
	KeyF12 = keycodes.TypeNormal | 0x45

	KeyPrintScreen = keycodes.TypeNormal | 0x46
	KeyScrollLock  = keycodes.TypeNormal | 0x47
	KeyPause       = keycodes.TypeNormal | 0x48
	KeyInsert      = keycodes.TypeNormal | 0x49
	KeyHome        = keycodes.TypeNormal | 0x4A
	KeyPageUp      = keycodes.TypeNormal | 0x4B
	KeyDelete      = keycodes.TypeNormal | 0x4C
	KeyEnd         = keycodes.TypeNormal | 0x4D
	KeyPageDown    = keycodes.TypeNormal | 0x4E
	KeyRight       = keycodes.TypeNormal | 0x4F
	KeyLeft        = keycodes.TypeNormal | 0x50
	KeyDown        = keycodes.TypeNormal | 0x51
	KeyUp          = keycodes.TypeNormal | 0x52

	KeyLeftCtrl   = keycodes.TypeNormal | 0xE0
	KeyLeftShift  = keycodes.TypeNormal | 0xE1
	KeyLeftAlt    = keycodes.TypeNormal | 0xE2
	KeyWindows    = keycodes.TypeNormal | 0xE3
	KeyRightCtrl  = keycodes.TypeNormal | 0xE4
	KeyRightShift = keycodes.TypeNormal | 0xE5
	KeyRightAlt   = keycodes.TypeNormal | 0xE6
	KeyRightGUI   = keycodes.TypeNormal | 0xE7
)

const (
	kcNO = 0
	// QMK KC_MEH and LCA(KC_LGUI). The low byte of KC_MEH is KC_NO.
	kcMEH  = keycodes.TypeXCtl | keycodes.TypeXSft | keycodes.TypeXAlt
	kcLCAG = keycodes.TypeXCtl | keycodes.TypeXAlt | 0xE3
)
