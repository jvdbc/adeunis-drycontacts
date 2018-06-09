package main

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	type args struct {
		ctx    context.Context
		sigfox sigfoxJSON
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{context.WithValue(nil, "test", testPost{}), sigfoxJSON{
			Time:      "2018-06-08T16:00:00.000Z",
			Device:    "2D4114",
			Duplicate: "false",
			Snr:       "14.29",
			Rssi:      "108.00",
			AvgSignal: "21.53",
			Station:   "station",
			Data:      "40ab00f100020001000001",
			Lat:       "45.751",
			Lng:       "4.860",
			SeqNumber: "0001",
		}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleRequest(tt.args.ctx, tt.args.sigfox)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HandleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testPost struct {
}

func (t testPost) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return &http.Response{}, nil
}
