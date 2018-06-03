package main

import (
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

// FrameCounter value from 3 bits (5-6-7)
func (us UplinkStatus) FrameCounter() uint8 {
	fc := us >> 5
	return uint8(fc)
}

// CmdOutputDone value from 1 bit (3)
func (us UplinkStatus) CmdOutputDone() bool {
	var mask UplinkStatus = 0x08 // 0000 1000
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
}

// UplinkFrame interface
type UplinkFrame interface {
}

// parseDevice func
func parseDevice(payload []byte) (UplinkFrame, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// parseNetwork func
func parseNetwork(payload []byte) (UplinkFrame, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// parseKeepalive func
func parseKeepalive(payload []byte) (UplinkFrame, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// parseResponse func
func parseResponse(payload []byte) (UplinkFrame, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// parseData func
func parseData(payload []byte) (UplinkFrame, error) {
	return nil, fmt.Errorf("Not Implemented")
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
