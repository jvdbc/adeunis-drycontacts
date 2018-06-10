package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func cleanEnv() {
	os.Unsetenv(callbackHost)
	os.Unsetenv(callbackPath)
}

func TestHandleRequest(t *testing.T) {
	defer cleanEnv()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	os.Setenv(callbackHost, url.Host)
	os.Setenv(callbackPath, url.Path)

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
		{"1", args{context.Background(), sigfoxJSON{
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
