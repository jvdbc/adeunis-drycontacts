package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {

}

// UplinkCode is the byte for possible uplink frame type
type UplinkCode byte

const (
	// Device config frame
	Device UplinkCode = 0x10

	// Network config frame
	Network UplinkCode = 0x20

	// Keepalive is a life frame
	Keepalive UplinkCode = 0x30

	// Response frame to a request or a config
	Response UplinkCode = 0x31

	// Data frame
	Data UplinkCode = 0x40
)

//go:generate stringer -type=UplinkCode

// UplinkStatus is the byte to reflect current device state
type UplinkStatus byte

// FrameCounter value from 3 bits (5-6-7) (from 0 value to 7 and restart)
func (us UplinkStatus) FrameCounter() uint8 {
	fc := us >> 5
	return uint8(fc)
}

// CmdOutputDone value from 1 bit (3)
func (us UplinkStatus) CmdOutputDone() bool {
	var mask UplinkStatus = 0x08 //0000 1000
	return ((us & mask) >> 3) > 0
}

// HWError value from 1 bit (2)
func (us UplinkStatus) HWError() bool {
	var mask UplinkStatus = 0x04 // 0000 0100
	return ((us & mask) >> 2) > 0
}

// LowBattery value from 1 bit (1)
func (us UplinkStatus) LowBattery() bool {
	var mask UplinkStatus = 0x02 // 0000 0010
	return ((us & mask) >> 1) > 0
}

// LastReqConfig value from 1 bit (0)
func (us UplinkStatus) LastReqConfig() bool {
	var mask UplinkStatus = 0x01 // 0000 0001
	return (us & mask) > 0
}

// Payload type to use Parse func
type Payload []byte

// Header 11 bytes for 4 tor data
type Header struct {
	Code   UplinkCode
	Status UplinkStatus
}

// DeviceFrame type
type DeviceFrame struct {
	*Header
}

// NetworkFrame type
type NetworkFrame struct {
	*Header
}

// KeepaliveFrame type
type KeepaliveFrame struct {
	*Header
}

// ResponseFrame type
type ResponseFrame struct {
	*Header
}

// DataFrame type
type DataFrame struct {
	*Header

	Tor1         uint16
	Tor1State    bool
	Tor1Previous bool

	Tor2         uint16
	Tor2State    bool
	Tor2Previous bool

	Tor3         uint16
	Tor3State    bool
	Tor3Previous bool

	Tor4         uint16
	Tor4State    bool
	Tor4Previous bool
}

// UplinkFrame interface
type UplinkFrame interface {
}

// parseDevice func
func parseDevice(payload []byte) (DeviceFrame, error) {
	return DeviceFrame{}, fmt.Errorf("Not Implemented")
}

// parseNetwork func
func parseNetwork(payload []byte) (NetworkFrame, error) {
	return NetworkFrame{}, fmt.Errorf("Not Implemented")
}

// parseKeepalive func
func parseKeepalive(payload []byte) (KeepaliveFrame, error) {
	return KeepaliveFrame{}, fmt.Errorf("Not Implemented")
}

// parseResponse func
func parseResponse(payload []byte) (ResponseFrame, error) {
	return ResponseFrame{}, fmt.Errorf("Not Implemented")
}

// parseData func
func parseData(payload []byte) (DataFrame, error) {
	frame := DataFrame{}

	if err := byteToUint16(payload[2:4], &frame.Tor1); err != nil {
		return frame, err
	}

	if err := byteToUint16(payload[4:6], &frame.Tor2); err != nil {
		return frame, err
	}

	if err := byteToUint16(payload[6:8], &frame.Tor3); err != nil {
		return frame, err
	}

	if err := byteToUint16(payload[8:10], &frame.Tor4); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 0, &frame.Tor1State); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 1, &frame.Tor1Previous); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 2, &frame.Tor2State); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 3, &frame.Tor2Previous); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 4, &frame.Tor3State); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 5, &frame.Tor3Previous); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 6, &frame.Tor4State); err != nil {
		return frame, err
	}

	if err := oneOrZero(payload[10], 7, &frame.Tor4Previous); err != nil {
		return frame, err
	}

	return frame, nil
}

func byteToUint16(slice []byte, result *uint16) error {
	err := binary.Read(bytes.NewReader(slice), binary.BigEndian, result)
	return err
}

func oneOrZero(value byte, bPos uint8, result *bool) error {
	if bPos > 7 {
		return fmt.Errorf("bit position should be between 0 and 8")
	}

	mask := byte(1 << bPos) // = 2 ^ bPosr (get mask)

	*result = ((value & mask) >> bPos) > 0

	return nil
}

// Parse func
func (p Payload) Parse() (UplinkFrame, error) {
	if len(p) != 11 {
		return nil, fmt.Errorf("Payload should have a size of 11 bytes")
	}

	switch UplinkCode(p[0]) {
	case Device:
		return parseDevice(p)
	case Network:
		return parseNetwork(p)
	case Keepalive:
		return parseKeepalive(p)
	case Response:
		return parseResponse(p)
	case Data:
		return parseData(p)
	default:
		return nil, fmt.Errorf("Unknown code byte")
	}
}
