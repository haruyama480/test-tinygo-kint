package main

import (
	"machine"
	"machine/usb/adc/midi"
	"time"
)

const (
	cable    = 0
	channel  = 1
	velocity = 0x40
)

var rows = []machine.Pin{
	machine.D8, machine.D9, machine.D10, machine.D11, machine.D7,
	machine.D16, machine.D5, machine.D3, machine.D4, machine.D1,
	machine.D0, machine.D2, machine.D17, machine.D23, machine.D21,
}

var cols = []machine.Pin{
	machine.D18, machine.D14, machine.D15,
	machine.D20, machine.D22, machine.D19, machine.D6,
}

func main() {
	m := midi.Port()

	for _, p := range rows {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.High()
	}
	for _, p := range cols {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	time.Sleep(2 * time.Second)

	var prev [15][7]bool
	for {
		for r, row := range rows {
			row.Low()
			time.Sleep(10 * time.Microsecond)
			for c, col := range cols {
				down := !col.Get()
				if down == prev[r][c] {
					continue
				}
				note := midi.Note(48 + r*7 + c) // C3 から割り当て
				if down {
					m.NoteOn(cable, channel, note, velocity)
				} else {
					m.NoteOff(cable, channel, note, velocity)
				}
				prev[r][c] = down
			}
			row.High()
		}
		time.Sleep(5 * time.Millisecond)
	}
}
