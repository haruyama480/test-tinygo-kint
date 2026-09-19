package main

type Keycode uint16

const (
	NumRows   = 15
	NumCols   = 7
	NumLayers = 2
)

// HID Keyboard/Keypad page (0x07). Names follow QMK KC_*.
const (
	KC_NO Keycode = 0x0000

	KC_A Keycode = 0x04
	KC_B Keycode = 0x05
	KC_C Keycode = 0x06
	KC_D Keycode = 0x07
	KC_E Keycode = 0x08
	KC_F Keycode = 0x09
	KC_G Keycode = 0x0A
	KC_H Keycode = 0x0B
	KC_I Keycode = 0x0C
	KC_J Keycode = 0x0D
	KC_K Keycode = 0x0E
	KC_L Keycode = 0x0F
	KC_M Keycode = 0x10
	KC_N Keycode = 0x11
	KC_O Keycode = 0x12
	KC_P Keycode = 0x13
	KC_Q Keycode = 0x14
	KC_R Keycode = 0x15
	KC_S Keycode = 0x16
	KC_T Keycode = 0x17
	KC_U Keycode = 0x18
	KC_V Keycode = 0x19
	KC_W Keycode = 0x1A
	KC_X Keycode = 0x1B
	KC_Y Keycode = 0x1C
	KC_Z Keycode = 0x1D

	KC_1 Keycode = 0x1E
	KC_2 Keycode = 0x1F
	KC_3 Keycode = 0x20
	KC_4 Keycode = 0x21
	KC_5 Keycode = 0x22
	KC_6 Keycode = 0x23
	KC_7 Keycode = 0x24
	KC_8 Keycode = 0x25
	KC_9 Keycode = 0x26
	KC_0 Keycode = 0x27

	KC_ENTER Keycode = 0x28
	KC_ESC   Keycode = 0x29
	KC_BSPC  Keycode = 0x2A
	KC_TAB   Keycode = 0x2B
	KC_SPC   Keycode = 0x2C
	KC_MINS  Keycode = 0x2D
	KC_EQL   Keycode = 0x2E
	KC_LBRC  Keycode = 0x2F
	KC_RBRC  Keycode = 0x30
	KC_BSLS  Keycode = 0x31
	KC_SCLN  Keycode = 0x33
	KC_QUOT  Keycode = 0x34
	KC_GRV   Keycode = 0x35
	KC_COMM  Keycode = 0x36
	KC_DOT   Keycode = 0x37
	KC_SLSH  Keycode = 0x38
	KC_CAPS  Keycode = 0x39

	KC_F1  Keycode = 0x3A
	KC_F2  Keycode = 0x3B
	KC_F3  Keycode = 0x3C
	KC_F4  Keycode = 0x3D
	KC_F5  Keycode = 0x3E
	KC_F6  Keycode = 0x3F
	KC_F7  Keycode = 0x40
	KC_F8  Keycode = 0x41
	KC_F9  Keycode = 0x42
	KC_F10 Keycode = 0x43
	KC_F11 Keycode = 0x44
	KC_F12 Keycode = 0x45

	KC_PSCR Keycode = 0x46
	KC_SCRL Keycode = 0x47
	KC_PAUS Keycode = 0x48
	KC_INS  Keycode = 0x49
	KC_HOME Keycode = 0x4A
	KC_PGUP Keycode = 0x4B
	KC_DEL  Keycode = 0x4C
	KC_END  Keycode = 0x4D
	KC_PGDN Keycode = 0x4E
	KC_RGHT Keycode = 0x4F
	KC_LEFT Keycode = 0x50
	KC_DOWN Keycode = 0x51
	KC_UP   Keycode = 0x52

	KC_LCTL Keycode = 0xE0
	KC_LSFT Keycode = 0xE1
	KC_LALT Keycode = 0xE2
	KC_LGUI Keycode = 0xE3
	KC_RCTL Keycode = 0xE4
	KC_RSFT Keycode = 0xE5
	KC_RALT Keycode = 0xE6
	KC_RGUI Keycode = 0xE7
)

// Boot report modifier bits (byte 0).
const (
	modLCtl = 1 << 0
	modLSft = 1 << 1
	modLAlt = 1 << 2
	modLGUI = 1 << 3
	modRCtl = 1 << 4
	modRSft = 1 << 5
	modRAlt = 1 << 6
	modRGUI = 1 << 7
)

// Firmware actions. Not HID usages; packBoot never emits these.
const (
	KC_BOOTLOADER Keycode = 0xF000 // QK_BOOT → machine.EnterBootloader
	KC_MEH        Keycode = 0xF001 // LCtl+LSft+LAlt
	KC_LCAG       Keycode = 0xF002 // LCtl+LAlt+LGUI (QMK LCA(KC_LGUI))
	KC_HYPR       Keycode = 0xF003 // LCtl+LSft+LAlt+LGUI

	// QMK dynamic macros. Goal 1 does not record/play; keys are no-ops.
	KC_DM_REC1 Keycode = 0xF010
	KC_DM_REC2 Keycode = 0xF011
	KC_DM_PLY1 Keycode = 0xF012
	KC_DM_PLY2 Keycode = 0xF013
	KC_DM_RSTP Keycode = 0xF014
)

func idx(row, col int) int { return row*NumCols + col }

func actionMods(kc Keycode) uint8 {
	switch kc {
	case KC_MEH:
		return modLCtl | modLSft | modLAlt
	case KC_LCAG:
		return modLCtl | modLAlt | modLGUI
	case KC_HYPR:
		return modLCtl | modLSft | modLAlt | modLGUI
	default:
		return 0
	}
}
