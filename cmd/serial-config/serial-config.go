package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tarm/serial"
)

func sendATOCommand(ser *serial.Port) bool {
	println("Send ATO")
	ser.Write([]byte("ATO\r\n"))
	reader := bufio.NewReader(ser)
	line := readLine(reader)
	if line == "O" {
		return true
	}
	println("ERROR")
	return false
}

func sendWriteCommand(ser *serial.Port, register string, value string) bool {
	println("Send Write command")
	ser.Write([]byte("ATS" + register + "=" + value + "\r\n"))
	reader := bufio.NewReader(ser)
	line := readLine(reader)
	println(line)
	if line == "O" {
		return true
	}
	println("ERROR")
	return false
}

func sendStoreCommand(ser *serial.Port) bool {
	println("Send AT &W")
	ser.Write([]byte("AT&W\r\n"))
	reader := bufio.NewReader(ser)
	line := readLine(reader)
	if line == "O" {
		return true
	}
	println("ERROR")
	return false
}

func initSequence(ser *serial.Port) bool {
	println("Send Init sequence")
	ser.Write([]byte("+++"))
	reader := bufio.NewReader(ser)
	line := readLine(reader)
	if line != "CONNECTING..." {
		if sendATOCommand(ser) {
			return true
		}
		println("ERROR")
		return false
	}
	line = readLine(reader)
	println(line)
	if line == "CM" {
		println("Init sequence OK")
		return true
	}
	println("ERROR")
	return false
}

func configurationSequence(ser *serial.Port) bool {
	println("Start configuration sequence")

	if !sendWriteCommand(ser, "300", "254") {
		return false
	}
	if !sendWriteCommand(ser, "301", "255") {
		return false
	}
	if !sendWriteCommand(ser, "302", "46") {
		return false
	}
	if !sendWriteCommand(ser, "303", "46") {
		return false
	}
	if !sendWriteCommand(ser, "304", "46") {
		return false
	}
	if !sendWriteCommand(ser, "305", "0") {
		return false
	}
	if !sendWriteCommand(ser, "307", "0") {
		return false
	}
	if !sendStoreCommand(ser) {
		return false
	}
	if !sendATOCommand(ser) {
		return false
	}
	return true
}

func close(ser *serial.Port) {
	if err := ser.Close(); err != nil {
		log.Fatal(err)
	}
}

// https://stackoverflow.com/questions/17599232/reading-from-serial-port-with-while-loop
func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [COM?] or [/dev/ttyUSB?]", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(-1)
	}

	com := strings.TrimSpace(os.Args[1])

	c := &serial.Config{Name: com, Baud: 115200, Parity: serial.ParityNone, StopBits: serial.Stop1, Size: 8}
	s, err := serial.OpenPort(c)

	if err != nil {
		log.Fatalf("Open Port %s : %s", com, err)
	}

	defer close(s)

	if initSequence(s) {
		if configurationSequence(s) {
			println("Configuraton OK")
		} else {
			println("Configuraton KO")
		}
	} else {
		println("Unable to configure")
	}

	println("End")
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(line)
}
