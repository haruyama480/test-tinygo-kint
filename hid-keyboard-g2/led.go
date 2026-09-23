package main

import "machine"

var (
	ledCaps    = machine.D12
	ledNum     = machine.D26
	ledScroll  = machine.D25
	ledCompose = machine.D24
	ledPower   = machine.D13
)

func configureLEDs() {
	for _, p := range []machine.Pin{ledCaps, ledNum, ledScroll, ledCompose} {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.High() // active-low, off
	}
	ledPower.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledPower.High()
}

func applyLEDs(v uint8) {
	hidLEDs = v
	writeActiveLow(ledNum, v&0x01 != 0)
	writeActiveLow(ledCaps, v&0x02 != 0)
	writeActiveLow(ledScroll, v&0x04 != 0)
	writeActiveLow(ledCompose, v&0x08 != 0)
}

func writeActiveLow(p machine.Pin, on bool) {
	if on {
		p.Low()
	} else {
		p.High()
	}
}
