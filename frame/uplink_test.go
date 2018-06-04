package frame

import (
	"reflect"
	"testing"
)

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

func TestPayload_Parse(t *testing.T) {
	tests := []struct {
		name    string
		p       Payload
		want    UplinkFrame
		wantErr bool
	}{
		{"1", Payload([]byte{}), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Payload.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Payload.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byteToUint16(t *testing.T) {
	type args struct {
		slice []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{"1", args{[]byte{0x00, 0x00}}, 0, false},
		{"2", args{[]byte{0x00, 0xFF}}, 255, false},
		{"3", args{[]byte{0x00, 0x01}}, 1, false},
		{"4", args{[]byte{0x00, 0x0A}}, 10, false},
		{"5", args{[]byte{0xFF, 0xFF}}, 65535, false},
		{"6", args{[]byte{0xFF, 0x00}}, 65280, false},
		{"7", args{[]byte{0x0A, 0x00}}, 2560, false},
		{"8", args{[]byte{0x0A, 0x0A}}, 2570, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := byteToUint16(tt.args.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("byteToUint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("byteToUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oneOrZero(t *testing.T) {
	type args struct {
		value byte
		bPos  uint8
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"1", args{0x00, 0}, false, false},
		{"2", args{0xFF, 0}, true, false},
		{"3", args{0xFF, 1}, true, false},
		{"4", args{0xFF, 2}, true, false},
		{"5", args{0xFF, 3}, true, false},
		{"6", args{0xFF, 4}, true, false},
		{"7", args{0xFF, 5}, true, false},
		{"8", args{0xFF, 6}, true, false},
		{"9", args{0xFF, 7}, true, false},
		{"10", args{0xFF, 8}, false, true},
		{"11", args{0x00, 9}, false, true},
		{"11", args{0x08, 0}, false, false},
		{"12", args{0x08, 1}, false, false},
		{"13", args{0x08, 2}, false, false},
		{"14", args{0x08, 3}, true, false},
		{"15", args{0x08, 4}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oneOrZero(tt.args.value, tt.args.bPos)
			if (err != nil) != tt.wantErr {
				t.Errorf("oneOrZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("oneOrZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseData(t *testing.T) {
	type args struct {
		header  *Header
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    DataFrame
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseData(tt.args.header, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseData() = %v, want %v", got, tt.want)
			}
		})
	}
}
