package main

// layer0 is QMK kinesis/keymaps/default_pretty (my-customize), stored as
// row*7+col. Matrix coordinates come from kint41/keyboard.json LAYOUT order.
var layer0 = [NumRows * NumCols]Keycode{
	// row 0
	KC_EQL, KC_TAB, KC_LCTL, KC_LSFT, KC_NO, KC_NO, KC_NO,
	// row 1
	KC_1, KC_Q, KC_A, KC_Z, KC_GRV, KC_NO, KC_NO,
	// row 2
	KC_2, KC_W, KC_S, KC_X, KC_INS, KC_LALT, KC_NO,
	// row 3
	KC_3, KC_E, KC_D, KC_C, KC_LEFT, KC_ESC, KC_SPC,
	// row 4
	KC_4, KC_R, KC_F, KC_V, KC_NO, KC_NO, KC_LGUI,
	// row 5
	KC_5, KC_T, KC_G, KC_B, KC_RGHT, KC_LCAG, KC_LCTL,
	// row 6
	KC_6, KC_Y, KC_H, KC_N, KC_DOWN, KC_ENTER, KC_MEH,
	// row 7
	KC_7, KC_U, KC_J, KC_M, KC_NO, KC_RGUI, KC_NO,
	// row 8
	KC_8, KC_I, KC_K, KC_COMM, KC_UP, KC_PGUP, KC_BSPC,
	// row 9
	KC_9, KC_O, KC_L, KC_DOT, KC_LBRC, KC_NO, KC_PGDN,
	// row 10
	KC_0, KC_P, KC_SCLN, KC_SLSH, KC_RBRC, KC_NO, KC_NO,
	// row 11
	KC_MINS, KC_BSLS, KC_QUOT, KC_RSFT, KC_NO, KC_NO, KC_NO,
	// row 12
	KC_ESC, KC_F3, KC_F6, KC_F9, KC_F12, KC_DM_PLY1, KC_NO,
	// row 13
	KC_F1, KC_F4, KC_F7, KC_F10, KC_DM_REC1, KC_DM_RSTP, KC_NO,
	// row 14
	KC_F2, KC_F5, KC_F8, KC_F11, KC_DM_REC2, KC_DM_PLY2, KC_NO,
}

var layers [NumLayers][NumRows * NumCols]Keycode

func init() {
	copy(layers[0][:], layer0[:])
}
