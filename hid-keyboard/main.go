package main

import (
	"machine"
	"time"
)

// debug prints press/release coordinates on CDC. Keep false for Goal 1 release.
const debug = false

var matrix Matrix

func main() {
	configureLEDs()
	matrix.configure()
	println("kinT TinyGo HID boot keyboard 16C0:0483 report", hidReportLen)

	var (
		pressed   [32]byte
		idleTicks uint32
	)
	period := time.Millisecond

	for {
		t0 := time.Now()
		matrix.scanOnce()
		mods, boot := collectKeys(&pressed)
		if boot {
			hidSendEmpty()
			machine.EnterBootloader()
		}
		report := packReport(&pressed, mods, hidProtocol)
		hidSend(report, &idleTicks)

		elapsed := time.Since(t0)
		if elapsed < period {
			time.Sleep(period - elapsed)
		}
	}
}

func collectKeys(pressed *[32]byte) (mods uint8, boot bool) {
	for i := range pressed {
		pressed[i] = 0
	}
	for r := 0; r < NumRows; r++ {
		for c := 0; c < NumCols; c++ {
			i := idx(r, c)
			down := matrix.stable[i]
			kc := layers[0][i]
			if down && !matrix.prev[i] {
				if debug {
					println("press", r, c, i, uint16(kc))
				}
				if kc == KC_BOOTLOADER {
					boot = true
				}
			} else if !down && matrix.prev[i] && debug {
				println("release", r, c, i, uint16(kc))
			}
			matrix.prev[i] = down
			if !down || kc == KC_NO || kc >= KC_BOOTLOADER {
				continue
			}
			if kc >= 0xE0 && kc <= 0xE7 {
				mods |= 1 << uint8(kc-0xE0)
				continue
			}
			if kc < 0xE8 {
				u := uint8(kc)
				pressed[u/8] |= 1 << (u % 8)
			}
		}
	}
	return mods, boot
}
