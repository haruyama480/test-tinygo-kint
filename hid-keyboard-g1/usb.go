package main

import (
	"machine"
	"machine/usb"
	"machine/usb/descriptor"
)

func init() {
	usb.Manufacturer = "kinT"
	usb.Product = "kinT TinyGo"

	machine.ConfigureUSBEndpoint(buildDescriptor(),
		[]usb.EndpointConfig{
			{
				Index:     usb.HID_ENDPOINT_OUT,
				IsIn:      false,
				Type:      usb.ENDPOINT_TYPE_INTERRUPT,
				RxHandler: hidRxLEDs,
			},
			{
				Index: usb.HID_ENDPOINT_IN,
				IsIn:  true,
				Type:  usb.ENDPOINT_TYPE_INTERRUPT,
			},
		},
		[]usb.SetupConfig{
			{Index: usb.HID_INTERFACE, Handler: hidSetup},
		},
	)
}

func buildDescriptor() descriptor.Descriptor {
	iface := append([]byte(nil), descriptor.InterfaceHID.Bytes()...)
	iface[6] = 0x01 // bInterfaceSubClass: Boot
	iface[7] = 0x01 // bInterfaceProtocol: Keyboard

	report := bootKeyboardReportDescriptor()
	hidReportLen = len(report)

	conf := descriptor.Append([][]byte{
		descriptor.ConfigurationCDCHID.Bytes(),
		descriptor.InterfaceAssociationCDC.Bytes(),
		descriptor.InterfaceCDCControl.Bytes(),
		descriptor.ClassSpecificCDCHeader.Bytes(),
		descriptor.ClassSpecificCDCACM.Bytes(),
		descriptor.ClassSpecificCDCUnion.Bytes(),
		descriptor.ClassSpecificCDCCallManagement.Bytes(),
		descriptor.EndpointIN(descriptor.EndpointEP1, descriptor.TransferTypeInterrupt, 0x10, 0x10).Bytes(),
		descriptor.InterfaceCDCData.Bytes(),
		descriptor.EndpointOUT(descriptor.EndpointEP2, descriptor.TransferTypeBulk, 0x40, 0x00).Bytes(),
		descriptor.EndpointIN(descriptor.EndpointEP2, descriptor.TransferTypeBulk, 0x40, 0x00).Bytes(),
		iface,
		descriptor.ClassHID.Bytes(),
		descriptor.EndpointIN(descriptor.EndpointEP3, descriptor.TransferTypeInterrupt, 8, 1).Bytes(),
		descriptor.EndpointOUT(descriptor.EndpointEP3, descriptor.TransferTypeInterrupt, 8, 1).Bytes(),
	})
	h, err := descriptor.FindClassHIDType(conf, descriptor.ClassHID.Bytes())
	if err != nil {
		panic(err)
	}
	h.ClassLength(uint16(len(report)))

	return descriptor.Descriptor{
		Device:        descriptor.DeviceCDC.Bytes(),
		Configuration: conf,
		HID:           map[uint16][]byte{usb.HID_INTERFACE: report},
	}
}

func bootKeyboardReportDescriptor() []byte {
	return descriptor.Append([][]byte{
		descriptor.HIDUsagePageGenericDesktop,
		descriptor.HIDUsageDesktopKeyboard,
		descriptor.HIDCollectionApplication,

		descriptor.HIDUsagePageKeyboard,
		descriptor.HIDUsageMinimum(224),
		descriptor.HIDUsageMaximum(231),
		descriptor.HIDLogicalMinimum(0),
		descriptor.HIDLogicalMaximum(1),
		descriptor.HIDReportSize(1),
		descriptor.HIDReportCount(8),
		descriptor.HIDInputDataVarAbs,

		descriptor.HIDReportCount(1),
		descriptor.HIDReportSize(8),
		descriptor.HIDInputConstVarAbs,

		descriptor.HIDUsagePageLED,
		descriptor.HIDUsageMinimum(1),
		descriptor.HIDUsageMaximum(5),
		descriptor.HIDReportCount(5),
		descriptor.HIDReportSize(1),
		descriptor.HIDOutputDataVarAbs,
		descriptor.HIDReportCount(1),
		descriptor.HIDReportSize(3),
		descriptor.HIDOutputConstVarAbs,

		descriptor.HIDUsagePageKeyboard,
		descriptor.HIDUsageMinimum(0),
		descriptor.HIDUsageMaximum(255),
		descriptor.HIDLogicalMinimum(0),
		descriptor.HIDLogicalMaximum(255),
		descriptor.HIDReportSize(8),
		descriptor.HIDReportCount(6),
		descriptor.HIDInputDataAryAbs,

		descriptor.HIDCollectionEnd,
	})
}
