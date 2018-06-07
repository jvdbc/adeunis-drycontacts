package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jvdbc/adeunis-drycontacts/frame"
)

// SigfoxEvent type
type SigfoxEvent struct {
	Name      string    `json:"name"`
	Time      time.Time `json:"time"`
	Device    string    `json:"device"`
	Duplicate string    `json:"duplicate"`
	Snr       string    `json:"snr"`
	Rssi      string    `json:"rssi"`
	AvgSignal string    `json:"avgSignal"`
	Station   string    `json:"station"`
	Data      string    `json:"data"`
	Lat       string    `json:"lat"`
	Lng       string    `json:"lng"`
	SeqNumber string    `json:"seqNumber"`
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
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Values    []IPTor   `json:"values"`
}

func init() {
}

// HandleRequest start point
func HandleRequest(ctx context.Context, sigfox SigfoxEvent) (string, error) {
	log.Printf("Received event: %v\n", sigfox)

	response(sigfox)

	return fmt.Sprintf("Hello %s!\n", sigfox.Name), nil
}

func response(sigfox SigfoxEvent) error {
	host := os.Getenv("smartconnectcallbackhost") // 'connector-demoenv.devinno.fr'
	path := os.Getenv("smartconnectcallbackPath") // /ip/tor

	payload := frame.Payload(sigfox.Data)

	uf, err := payload.Parse()

	if err != nil {
		return fmt.Errorf("iptor parse error : %s", err)
	}

	var df frame.DataFrame

	switch uf.Code() {
	case frame.Data:
		df = uf.(frame.DataFrame)
	default:
		return fmt.Errorf("iptor frame not implemented: %s", uf.Code())
	}

	content := IPTorEvent{sigfox.Device, sigfox.Time, []IPTor{
		IPTor{1, "", true, df.Tor1 > 0},
		IPTor{2, "", true, df.Tor2 > 0},
		IPTor{3, "", true, df.Tor3 > 0},
		IPTor{4, "", true, df.Tor4 > 0},
	}}

	url := fmt.Sprintf("https://%v%v", host, path)
	body, err := json.Marshal(content)

	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("post response error: %v", err)
	}

	defer res.Body.Close()
	rspBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("body response error: %v", err)
	}

	log.Printf("Response: %s\n", string(rspBody))

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
