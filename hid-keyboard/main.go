package main

import (
	"machine"
	"machine/usb/hid/keyboard"
	"time"
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
	kb := keyboard.Port()
	machine.USBDev.Configure(machine.UARTConfig{})

	for _, p := range rows {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.High()
	}
	for _, p := range cols {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	time.Sleep(2 * time.Second)
	pressed := false

	for {
		down := scan()
		if down && !pressed {
			kb.Write([]byte("tinygo"))
		}
		pressed = down
		time.Sleep(5 * time.Millisecond)
	}
}

func scan() bool {
	for _, row := range rows {
		row.Low()
		time.Sleep(time.Microsecond * 10)
		for _, col := range cols {
			if !col.Get() {
				row.High()
				return true
			}
		}
		row.High()
	}
	return false
}
