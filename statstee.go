package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rodaine/statstee/datagram"
	"github.com/rodaine/statstee/router"
	"github.com/rodaine/statstee/views"
)

var (
	errBuffer bytes.Buffer

	deviceInterface string = "lo"
	sniffedPort     int    = 8125
	windowSize      int    = 600
	outputDebug     bool   = false
)

func init() {
	log.SetOutput(&errBuffer)

	flag.StringVar(&deviceInterface, "d", deviceInterface, "network device to listen on")
	flag.IntVar(&sniffedPort, "p", sniffedPort, "statsd UDP port to listen on")
	flag.IntVar(&windowSize, "n", windowSize, "seconds of data to keep")
	flag.BoolVar(&outputDebug, "v", outputDebug, "display any debug output on quit")
	flag.Parse()
}

func main() {
	if outputDebug {
		defer fmt.Fprintf(os.Stdout, errBuffer.String())
	}

	c := make(chan datagram.Metric, windowSize)
	router := router.New(c)

	go captureMetrics(router)
	go sniffStream(c)

	logIfError(views.Display(router))
}

func sniffStream(c chan datagram.Metric) {
	logIfError(datagram.Stream(deviceInterface, sniffedPort, c))
	close(c)
}

func captureMetrics(r *router.Router) {
	r.Listen()
	views.Quit()
}

func logIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}
