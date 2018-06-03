package main

import (
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestUplinkStatus_FrameCounter(t *testing.T) {
	tests := []struct {
		name string
		us   UplinkStatus
		want uint8
	}{
		{"1", 0xFF, 7},
		{"2", 0x00, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.FrameCounter(); got != tt.want {
				t.Errorf("UplinkStatus.FrameCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplinkStatus_CmdOutputDone(t *testing.T) {
	tests := []struct {
		name string
		us   UplinkStatus
		want bool
	}{
		{"1", 0x00, false},
		{"2", 0x08, true},
		{"3", 0x09, true},
		{"4", 0x10, false},
		{"5", 0x11, false},
		{"6", 0xFF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.CmdOutputDone(); got != tt.want {
				t.Errorf("UplinkStatus.CmdOutputDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplinkStatus_HWError(t *testing.T) {
	tests := []struct {
		name string
		us   UplinkStatus
		want bool
	}{
		{"1", 0x00, false},
		{"2", 0x01, false},
		{"3", 0x02, false},
		{"4", 0x03, false},
		{"5", 0x04, true},
		{"6", 0x05, true},
		{"8", 0x06, true},
		{"9", 0x07, true},
		{"10", 0x08, false},
		{"11", 0xFF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.HWError(); got != tt.want {
				t.Errorf("UplinkStatus.HWError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplinkStatus_LowBattery(t *testing.T) {
	tests := []struct {
		name string
		us   UplinkStatus
		want bool
	}{
		{"1", 0x00, false},
		{"2", 0x01, false},
		{"3", 0x02, true},
		{"4", 0x03, true},
		{"5", 0x04, false},
		{"6", 0x05, false},
		{"7", 0x06, true},
		{"8", 0xFF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.LowBattery(); got != tt.want {
				t.Errorf("UplinkStatus.LowBattery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUplinkStatus_LastReqConfig(t *testing.T) {
	tests := []struct {
		name string
		us   UplinkStatus
		want bool
	}{
		{"1", 0x00, false},
		{"2", 0x01, true},
		{"3", 0x02, false},
		{"4", 0x03, true},
		{"5", 0x04, false},
		{"6", 0x05, true},
		{"7", 0xFF, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.LastReqConfig(); got != tt.want {
				t.Errorf("UplinkStatus.LastReqConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
