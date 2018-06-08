package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jvdbc/adeunis-drycontacts/frame"
)

// SigfoxEvent type
type SigfoxEvent struct {
	Name      string `json:"name"`
	Time      string `json:"time"`
	Device    string `json:"device"`
	Duplicate string `json:"duplicate"`
	Snr       string `json:"snr"`
	Rssi      string `json:"rssi"`
	AvgSignal string `json:"avgSignal"`
	Station   string `json:"station"`
	Data      string `json:"data"`
	Lat       string `json:"lat"`
	Lng       string `json:"lng"`
	SeqNumber string `json:"seqNumber"`
}

// IPTor type
type IPTor struct {
	ID      uint8  `json:"id"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
	State   bool   `json:"state"`
}

// IPTorEvent type
type IPTorEvent struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Values    []IPTor `json:"values"`
}

func init() {
}

// HandleRequest start point
func HandleRequest(ctx context.Context, sigfox SigfoxEvent) (string, error) {
	log.Printf("Received event: %v\n", sigfox)

	payload, err := hex.DecodeString(sigfox.Data)

	if err != nil {
		return "", fmt.Errorf("hex decode fail: %v", err)
	}

	uf, err := frame.Payload(payload).Parse()

	if err != nil {
		return "", fmt.Errorf("iptor parse error: %v", err)
	}

	var df frame.DataFrame

	switch uf.Code() {
	case frame.Data:
		df = uf.(frame.DataFrame)
	default:
		return "", fmt.Errorf("iptor frame not implemented: %s", uf.Code())
	}

	content := IPTorEvent{sigfox.Device, sigfox.Time, []IPTor{
		IPTor{1, "alerte", true, df.Tor1State},
		IPTor{2, "alerte", true, df.Tor2State},
		IPTor{3, "alerte", true, df.Tor3State},
		IPTor{4, "alerte", true, df.Tor4State},
	}}

	host := os.Getenv("scCallbackHost") // 'connector-demoenv.devinno.fr'
	path := os.Getenv("scCallbackPath") // /ip/tor/data

	url := fmt.Sprintf("https://%v%v", host, path)
	log.Printf("call url: %v", url)

	if err != nil {
		return "", fmt.Errorf("json marshal error: %v", err)
	}

	log.Printf("json send: %+v", content)

	body, err := json.Marshal(content)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("post response error: %v", err)
	}

	defer res.Body.Close()
	rspBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", fmt.Errorf("body response error: %v", err)
	}

	log.Printf("Response: %v", string(rspBody))

	return fmt.Sprintf("Done !"), nil
}

func main() {
	lambda.Start(HandleRequest)
}
