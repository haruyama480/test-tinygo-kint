// Flash smoke test for machine.Flash. One boot erases the first block and
// writes a 3632-byte record, the same size tinygo-keyboard Save uses.
// The magic is not Flash.Size(), so a later keyboard boot ignores the record.
//
// PASS on this boot only shows that erase, program, and readback agree.
// Unplug and plug back in, then open the serial port again. The next boot
// must print the generation written here as previous-generation.
// Reflash this same program. The generation must still be there: the Teensy
// bootloader keeps this region.
package main

import (
	"machine"
	"machine/usb"
	"time"
)

// 4 byte header in Save, 6 layers, 105 keys, 2 bytes each, 2048 byte macro
// buffer, 32 combos of 5 keycodes.
const recordLen = 4 + 6*105*2 + 2048 + 32*5*2

var magic = [4]byte{'k', 'F', 'L', '1'}

func init() {
	usb.Product = "kint-flash-test"
}

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	time.Sleep(2 * time.Second)

	println("size", machine.Flash.Size())
	println("erase-block", machine.Flash.EraseBlockSize())
	println("write-block", machine.Flash.WriteBlockSize())
	println("data-start", uint32(machine.FlashDataStart()))
	println("data-end", uint32(machine.FlashDataEnd()))

	prev := make([]byte, recordLen)
	if _, err := machine.Flash.ReadAt(prev, 0); err != nil {
		fail(err)
	}
	prevGen, prevOK := parse(prev)
	if prevOK {
		println("previous-generation", prevGen)
	} else {
		println("previous-generation none")
	}

	nextGen := uint32(1)
	if prevOK {
		nextGen = prevGen + 1
	}
	rec := encode(nextGen)

	needed := int64(len(rec)) / machine.Flash.EraseBlockSize()
	if int64(len(rec))%machine.Flash.EraseBlockSize() != 0 {
		needed++
	}
	println("erasing", needed)
	t0 := time.Now()
	if err := machine.Flash.EraseBlocks(0, needed); err != nil {
		fail(err)
	}
	println("erase-ms", time.Since(t0).Milliseconds())

	println("writing", len(rec))
	t0 = time.Now()
	if _, err := machine.Flash.WriteAt(rec, 0); err != nil {
		fail(err)
	}
	println("write-ms", time.Since(t0).Milliseconds())

	got := make([]byte, len(rec))
	if _, err := machine.Flash.ReadAt(got, 0); err != nil {
		fail(err)
	}
	for i := range rec {
		if rec[i] != got[i] {
			println("FAIL mismatch", i)
			fail(nil)
		}
	}
	println("PASS", nextGen)

	for {
		if prevOK {
			println("PASS generation", nextGen, "previous", prevGen)
		} else {
			println("PASS generation", nextGen, "previous none")
		}
		for i := 0; i < 3; i++ {
			blink(true)
		}
	}
}

func encode(gen uint32) []byte {
	b := make([]byte, recordLen)
	copy(b[:4], magic[:])
	b[4] = byte(gen)
	b[5] = byte(gen >> 8)
	b[6] = byte(gen >> 16)
	b[7] = byte(gen >> 24)
	for i := 12; i < len(b); i++ {
		b[i] = byte(gen) ^ byte(i) ^ byte(i>>8)
	}
	sum := checksum(b[12:], gen)
	b[8] = byte(sum)
	b[9] = byte(sum >> 8)
	b[10] = byte(sum >> 16)
	b[11] = byte(sum >> 24)
	return b
}

func parse(b []byte) (uint32, bool) {
	if len(b) < recordLen {
		return 0, false
	}
	for i := range magic {
		if b[i] != magic[i] {
			return 0, false
		}
	}
	gen := uint32(b[4]) | uint32(b[5])<<8 | uint32(b[6])<<16 | uint32(b[7])<<24
	sum := uint32(b[8]) | uint32(b[9])<<8 | uint32(b[10])<<16 | uint32(b[11])<<24
	if sum != checksum(b[12:], gen) {
		return gen, false
	}
	return gen, true
}

func checksum(p []byte, gen uint32) uint32 {
	s := gen
	for _, c := range p {
		s += uint32(c)
	}
	return s
}

func fail(err error) {
	if err != nil {
		println("FAIL", err.Error())
	}
	for {
		blink(false)
	}
}

func blink(pass bool) {
	on := 200 * time.Millisecond
	off := 800 * time.Millisecond
	if !pass {
		on = 80 * time.Millisecond
		off = 80 * time.Millisecond
	}
	machine.LED.High()
	time.Sleep(on)
	machine.LED.Low()
	time.Sleep(off)
}
