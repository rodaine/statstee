package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rodaine/statstee/datagram"
	"github.com/rodaine/statstee/router"
	"github.com/rodaine/statstee/views"
)

const (
	logFile              = "statstee.log"
	metricBufferSize int = 600
)

var (
	logWriter   io.WriteCloser
	streamError error

	deviceInterface string = "lo"
	sniffedPort     int    = 8125
	outputDebug     bool   = false
)

func init() {
	flag.StringVar(&deviceInterface, "d", deviceInterface, "network device to listen on")
	flag.IntVar(&sniffedPort, "p", sniffedPort, "statsd UDP port to listen on")
	flag.BoolVar(&outputDebug, "v", outputDebug, "display debug output to "+logFile)
	flag.Parse()

	log.SetOutput(ioutil.Discard)
}

func main() {
	if outputDebug {
		logWriter, _ = os.OpenFile(logFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		log.SetFlags(log.Lmicroseconds | log.LstdFlags | log.Lshortfile)
		log.SetOutput(logWriter)
		defer logWriter.Close()
	}

	c := make(chan datagram.Metric, metricBufferSize)
	router := router.New(c)

	go captureMetrics(router)
	go sniffStream(c)

	fatalIfError(views.Loop(router))
	fatalIfError(streamError)
}

func sniffStream(c chan datagram.Metric) {
	streamError = datagram.Stream(deviceInterface, sniffedPort, c)
	close(c)
}

func captureMetrics(r *router.Router) {
	r.Listen()
	views.Quit()
}

func fatalIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
