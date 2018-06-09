package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
	code   UplinkCode
	status UplinkStatus
}

// Code get header code
func (h Header) Code() UplinkCode {
	return h.code
}

// Status get header status
func (h Header) Status() UplinkStatus {
	return h.status
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
	Code() UplinkCode
	Status() UplinkStatus
}

// Parse func
func (p Payload) Parse() (UplinkFrame, error) {
	if len(p) != 11 {
		return nil, fmt.Errorf("payload should have a size of 11 bytes, %s", p)
	}
	code := UplinkCode(p[0])
	header := &Header{code: code, status: UplinkStatus(p[1])}

	switch UplinkCode(code) {
	case Device:
		return parseDevice(header, p[2:])
	case Network:
		return parseNetwork(header, p[2:])
	case Keepalive:
		return parseKeepalive(header, p[2:])
	case Response:
		return parseResponse(header, p[2:])
	case Data:
		return parseData(header, p[2:])
	default:
		return nil, fmt.Errorf("unknown code byte")
	}
}

// parseDevice func
func parseDevice(header *Header, payload []byte) (DeviceFrame, error) {
	return DeviceFrame{}, fmt.Errorf("not implemented")
}

// parseNetwork func
func parseNetwork(header *Header, payload []byte) (NetworkFrame, error) {
	return NetworkFrame{}, fmt.Errorf("not implemented")
}

// parseKeepalive func
func parseKeepalive(header *Header, payload []byte) (KeepaliveFrame, error) {
	return KeepaliveFrame{}, fmt.Errorf("not implemented")
}

// parseResponse func
func parseResponse(header *Header, payload []byte) (ResponseFrame, error) {
	return ResponseFrame{}, fmt.Errorf("not implemented")
}

// parseData func
func parseData(header *Header, payload []byte) (DataFrame, error) {
	frame := DataFrame{Header: header}
	var err error

	if frame.Tor1, err = byteToUint16(payload[0:2]); err != nil {
		return frame, err
	}

	if frame.Tor2, err = byteToUint16(payload[2:4]); err != nil {
		return frame, err
	}

	if frame.Tor3, err = byteToUint16(payload[4:6]); err != nil {
		return frame, err
	}

	if frame.Tor4, err = byteToUint16(payload[6:8]); err != nil {
		return frame, err
	}

	if frame.Tor1State, err = oneOrZero(payload[8], 0); err != nil {
		return frame, err
	}

	if frame.Tor1Previous, err = oneOrZero(payload[8], 1); err != nil {
		return frame, err
	}

	if frame.Tor2State, err = oneOrZero(payload[8], 2); err != nil {
		return frame, err
	}

	if frame.Tor2Previous, err = oneOrZero(payload[8], 3); err != nil {
		return frame, err
	}

	if frame.Tor3State, err = oneOrZero(payload[8], 4); err != nil {
		return frame, err
	}

	if frame.Tor3Previous, err = oneOrZero(payload[8], 5); err != nil {
		return frame, err
	}

	if frame.Tor4State, err = oneOrZero(payload[8], 6); err != nil {
		return frame, err
	}

	if frame.Tor4Previous, err = oneOrZero(payload[8], 7); err != nil {
		return frame, err
	}

	return frame, nil
}

func byteToUint16(slice []byte) (uint16, error) {
	var r uint16
	err := binary.Read(bytes.NewReader(slice), binary.BigEndian, &r)
	return r, err
}

func oneOrZero(value byte, bPos uint8) (bool, error) {
	if bPos > 7 {
		return false, fmt.Errorf("bit position should be between 0 and 8")
	}
	mask := byte(1 << bPos) // = 2 ^ bPosr (get mask)

	return ((value & mask) >> bPos) > 0, nil
}
