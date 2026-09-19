package main

import (
	"context"
	"log"
	"machine"
	"machine/usb"

	keyboard "github.com/sago35/tinygo-keyboard"
	jp "github.com/sago35/tinygo-keyboard/keycodes/japanese"
)

func main() {
	usb.Product = "tinygo-kint-kb500"

	d := keyboard.New()

	colPins := []machine.Pin{
		machine.D18, machine.D14, machine.D15, machine.D20,
		machine.D22, machine.D19, machine.D6,
	}
	rowPins := []machine.Pin{
		machine.D8, machine.D9, machine.D10, machine.D11, machine.D7,
		machine.D16, machine.D5, machine.D3, machine.D4, machine.D1,
		machine.D0, machine.D2, machine.D17, machine.D23, machine.D21,
	}

	// layer0[row*7+col]
	layer0 := make([]keyboard.Keycode, 15*7)

	// 左手キーウェル
	layer0[0] = jp.KeyHat       // [0,0] =
	layer0[7] = jp.Key1         // [1,0]
	layer0[14] = jp.Key2        // [2,0]
	layer0[21] = jp.Key3        // [3,0]
	layer0[28] = jp.Key4        // [4,0]
	layer0[35] = jp.Key5        // [5,0]
	layer0[1] = jp.KeyTab       // [0,1]
	layer0[8] = jp.KeyQ         // [1,1]
	layer0[15] = jp.KeyW        // [2,1]
	layer0[22] = jp.KeyE        // [3,1]
	layer0[29] = jp.KeyR        // [4,1]
	layer0[36] = jp.KeyT        // [5,1]
	layer0[2] = jp.KeyCapsLock  // [0,2]
	layer0[9] = jp.KeyA         // [1,2]
	layer0[16] = jp.KeyS        // [2,2]
	layer0[23] = jp.KeyD        // [3,2]
	layer0[30] = jp.KeyF        // [4,2]
	layer0[37] = jp.KeyG        // [5,2]
	layer0[3] = jp.KeyLeftShift // [0,3]
	layer0[10] = jp.KeyZ        // [1,3]
	layer0[17] = jp.KeyX        // [2,3]
	layer0[24] = jp.KeyC        // [3,3]
	layer0[31] = jp.KeyV        // [4,3]
	layer0[38] = jp.KeyB        // [5,3]

	// 右手キーウェル
	layer0[42] = jp.Key6 // [6,0]
	layer0[49] = jp.Key7
	layer0[56] = jp.Key8
	layer0[63] = jp.Key9
	layer0[70] = jp.Key0
	layer0[77] = jp.KeyMinus // [11,0] -
	layer0[43] = jp.KeyY     // [6,1]
	layer0[50] = jp.KeyU
	layer0[57] = jp.KeyI
	layer0[64] = jp.KeyO
	layer0[71] = jp.KeyP
	layer0[78] = jp.KeyBackslash // [11,1]
	layer0[44] = jp.KeyH         // [6,2]
	layer0[51] = jp.KeyJ
	layer0[58] = jp.KeyK
	layer0[65] = jp.KeyL
	layer0[72] = jp.KeySemicolon
	layer0[79] = jp.KeyColon // US なら KeyQuote 相当。JP配列向け
	layer0[45] = jp.KeyN     // [6,3]
	layer0[52] = jp.KeyM
	layer0[59] = jp.KeyComma
	layer0[66] = jp.KeyPeriod
	layer0[73] = jp.KeySlash
	layer0[80] = jp.KeyRightShift // [11,3]

	// ファンクション行
	layer0[84] = jp.KeyEsc // [12,0]
	layer0[91] = jp.KeyF1  // [13,0]
	layer0[98] = jp.KeyF2  // [14,0]
	layer0[85] = jp.KeyF3  // [12,1]
	layer0[92] = jp.KeyF4
	layer0[99] = jp.KeyF5
	layer0[86] = jp.KeyF6 // [12,2]
	layer0[93] = jp.KeyF7
	layer0[100] = jp.KeyF8
	layer0[87] = jp.KeyF9 // [12,3]
	layer0[94] = jp.KeyF10
	layer0[101] = jp.KeyF11
	layer0[88] = jp.KeyF12 // [12,4]
	layer0[95] = jp.KeyPrintscreen
	layer0[102] = jp.KeyScrollLock
	layer0[89] = jp.KeyPause                // [12,5]
	layer0[103] = jp.KeyNumLock             // [14,5]
	layer0[96] = jp.KeyRestoreDefaultKeymap // [13,5] 初期化用

	// 下段の小さいキー
	layer0[11] = jp.KeyHankaku    // [1,4]
	layer0[18] = jp.KeyLeft       // [2,4]
	layer0[25] = jp.KeyRight      // [3,4]
	layer0[39] = jp.KeyHome       // [5,4]
	layer0[46] = jp.KeyUp         // [6,4]
	layer0[60] = jp.KeyDown       // [8,4]
	layer0[67] = jp.KeyLeftBrace  // [9,4]
	layer0[74] = jp.KeyRightBrace // [10,4]

	// 親指クラスタ
	layer0[27] = jp.KeyDelete    // [3,6] 左縦長
	layer0[34] = jp.KeyBackspace // [4,6] 左縦長
	layer0[19] = jp.KeyEnd       // [2,5]
	layer0[26] = jp.KeyHome      // [3,5]
	layer0[40] = jp.KeyLeftCtrl  // [5,5]
	layer0[41] = jp.KeyLeftAlt   // [5,6]
	layer0[47] = jp.KeySpace     // [6,5] 右縦長
	layer0[54] = jp.KeyEnter     // [7,5] 右縦長
	layer0[48] = jp.KeyPageDown  // [6,6]
	layer0[62] = jp.KeyPageUp    // [8,6]
	layer0[61] = jp.KeyRightAlt  // [8,5]
	layer0[69] = jp.KeyRightCtrl // [9,6]

	d.AddMatrixKeyboard(colPins, rowPins, [][]keyboard.Keycode{
		layer0,
	}, keyboard.InvertDiode(true))

	loadKeyboardDef()
	if err := d.Loop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
