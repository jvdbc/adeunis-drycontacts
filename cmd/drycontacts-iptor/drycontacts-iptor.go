package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jvdbc/adeunis-drycontacts/frame"
)

// sigfoxJSON type
type sigfoxJSON struct {
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

// iPTor type
type iPTor struct {
	ID      uint8  `json:"id"`
	Label   string `json:"label"`
	Enabled bool   `json:"enabled"`
	State   bool   `json:"state"`
}

// iPTorJSON type
type iPTorJSON struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Values    []iPTor `json:"values"`
}

type httpPoster interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type httpPost struct {
}

func (d httpPost) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}

var clt httpPoster

// HandleRequest start point
func HandleRequest(ctx context.Context, sigfox sigfoxJSON) (string, error) {

	log.Printf("received event: %v", sigfox)

	data, err := hex.DecodeString(sigfox.Data)
	if err != nil {
		return "", fmt.Errorf("hex decode failed: %v", err)
	}

	uf, err := frame.Payload(data).Parse()
	if err != nil {
		return "", fmt.Errorf("frame parse failed: %v", err)
	}

	var df frame.DataFrame
	switch x := uf.(type) {
	case frame.DataFrame:
		df = x
	default:
		return "", fmt.Errorf("%t frame not implemented", x)
	}

	content := iPTorJSON{
		sigfox.Device,
		sigfox.Time,
		[]iPTor{
			iPTor{1, "alerte", true, df.Tor1State},
			iPTor{2, "alerte", true, df.Tor2State},
			iPTor{3, "alerte", true, df.Tor3State},
			iPTor{4, "alerte", true, df.Tor4State},
		}}

	body, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %v", err)
	}

	log.Printf("json send: %+v", content)

	host := os.Getenv("scCallbackHost") // 'connector-demoenv.devinno.fr'
	path := os.Getenv("scCallbackPath") // /ip/tor/data

	url := fmt.Sprintf("https://%v%v", host, path)
	log.Printf("call url: %v", url)

	if ctx != nil {
		if tst, ok := ctx.Value("test").(httpPoster); ok {
			clt = tst
		}
	}

	if clt == nil {
		clt = httpPost{}
	}

	res, err := clt.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("post response failed: %v", err)
	}

	defer res.Body.Close()
	rspBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read body response failed: %v", err)
	}

	log.Printf("response: %v", string(rspBody))

	return fmt.Sprintf("lambda success"), nil
}

func main() {
	lambda.Start(HandleRequest)
}
