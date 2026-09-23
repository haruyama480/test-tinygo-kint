package main

import (
	"machine"
	"unsafe"

	"runtime/volatile"
)

const (
	debounceTicks   = 5
	selectSettleUs  = 5
	unselectDelayUs = 20
	cyclesPerUs     = 600 // Teensy 4.1 CORE_FREQ = 600 MHz
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

var dwtCYCCNT = (*volatile.Register32)(unsafe.Pointer(uintptr(0xe0001004)))

type Matrix struct {
	count  [NumRows * NumCols]uint8
	stable [NumRows * NumCols]bool
	prev   [NumRows * NumCols]bool
}

func delayUs(us uint32) {
	start := dwtCYCCNT.Get()
	cycles := us * cyclesPerUs
	for dwtCYCCNT.Get()-start < cycles {
	}
}

func (m *Matrix) configure() {
	for _, p := range rows {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.High()
	}
	for _, p := range cols {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
}

func (m *Matrix) scanOnce() {
	for r, row := range rows {
		row.Low()
		delayUs(selectSettleUs)
		for c, col := range cols {
			m.debounce(r, c, !col.Get())
		}
		row.High()
		delayUs(unselectDelayUs)
	}
}

func (m *Matrix) debounce(r, c int, raw bool) {
	i := idx(r, c)
	if raw == m.stable[i] {
		m.count[i] = 0
		return
	}
	m.count[i]++
	if m.count[i] >= debounceTicks {
		m.stable[i] = raw
		m.count[i] = 0
	}
}
