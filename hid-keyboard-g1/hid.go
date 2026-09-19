package main

import (
	"machine"
	"machine/usb"
	"runtime/interrupt"
)

const (
	protocolBoot   = 0
	protocolReport = 1
	reportTypeIn   = 1
	reportTypeOut  = 2
)

var (
	hidProtocol  uint8 = protocolReport
	hidIdle      uint8
	hidLEDs      uint8
	hidCurrent   [8]byte
	hidLastIN    [8]byte
	hidEP0       [8]byte
	hidSendFail  uint32
	hidReportLen int
)

func hidSetup(setup usb.Setup) bool {
	switch setup.BmRequestType {
	case usb.REQUEST_HOSTTODEVICE_CLASS_INTERFACE:
		switch setup.BRequest {
		case usb.SET_IDLE:
			hidIdle = setup.WValueH
			machine.SendZlp()
			return true
		case usb.SET_PROTOCOL:
			hidProtocol = setup.WValueL
			machine.SendZlp()
			return true
		case usb.SET_REPORT:
			if setup.WLength == 0 {
				machine.SendZlp()
				return true
			}
			if setup.WValueH != reportTypeOut || setup.WValueL != 0 {
				return false
			}
			b, err := machine.ReceiveUSBControlPacket()
			if err != nil {
				return false
			}
			applyLEDs(b[0])
			machine.SendZlp()
			return true
		}
	case usb.REQUEST_DEVICETOHOST_CLASS_INTERFACE:
		switch setup.BRequest {
		case usb.GET_IDLE:
			hidEP0[0] = hidIdle
			return machine.SendUSBInPacket(0, hidEP0[:1])
		case usb.GET_PROTOCOL:
			hidEP0[0] = hidProtocol
			return machine.SendUSBInPacket(0, hidEP0[:1])
		case usb.GET_REPORT:
			if setup.WValueL != 0 {
				return false
			}
			switch setup.WValueH {
			case reportTypeIn:
				var local [8]byte
				state := interrupt.Disable()
				local = hidLastIN
				interrupt.Restore(state)
				n := int(setup.WLength)
				if n > 8 {
					n = 8
				}
				copy(hidEP0[:], local[:])
				return machine.SendUSBInPacket(0, hidEP0[:n])
			case reportTypeOut:
				hidEP0[0] = hidLEDs
				return machine.SendUSBInPacket(0, hidEP0[:1])
			default:
				return false
			}
		}
	}
	return false
}

func hidRxLEDs(b []byte) {
	if len(b) == 0 {
		return
	}
	applyLEDs(b[0])
}

func packReport(pressed *[32]byte, mods uint8, protocol uint8) [8]byte {
	_ = protocol // Goal 1: always boot 8-byte
	return packBoot(pressed, mods)
}

func packBoot(pressed *[32]byte, mods uint8) (out [8]byte) {
	out[0] = mods
	n := 0
	for usage := uint16(0x04); usage <= 0xDF; usage++ {
		u := uint8(usage)
		if pressed[u/8]&(1<<(u%8)) == 0 {
			continue
		}
		n++
		if n > 6 {
			for i := 2; i < 8; i++ {
				out[i] = 0x01
			}
			return out
		}
		out[1+n] = u
	}
	return out
}

func hidSend(report [8]byte, idleTicks *uint32) {
	hidCurrent = report
	changed := report != hidLastIN
	idle := hidIdle
	dueIdle := idle != 0 && *idleTicks >= uint32(idle)*4
	if !changed && !dueIdle {
		*idleTicks++
		return
	}
	if !machine.SendUSBInPacket(uint32(usb.HID_ENDPOINT_IN), hidCurrent[:]) {
		hidSendFail++
		if hidSendFail == 1 || hidSendFail == 10 || hidSendFail == 100 || hidSendFail%1000 == 0 {
			println("hid in busy", hidSendFail)
		}
		*idleTicks++
		return
	}
	hidSendFail = 0
	*idleTicks = 0
	state := interrupt.Disable()
	hidLastIN = hidCurrent
	interrupt.Restore(state)
}

func hidSendEmpty() {
	var empty [8]byte
	_ = machine.SendUSBInPacket(uint32(usb.HID_ENDPOINT_IN), empty[:])
}
