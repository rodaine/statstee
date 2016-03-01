package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rodaine/statstee/datagram"
	"github.com/rodaine/statstee/router"
	"github.com/rodaine/statstee/views"
)

var (
	logWriter io.WriteCloser

	deviceInterface string = "lo"
	sniffedPort     int    = 8125
	windowSize      int    = 600
	outputDebug     bool   = false
)

func init() {
	flag.StringVar(&deviceInterface, "d", deviceInterface, "network device to listen on")
	flag.IntVar(&sniffedPort, "p", sniffedPort, "statsd UDP port to listen on")
	flag.IntVar(&windowSize, "n", windowSize, "seconds of data to keep")
	flag.BoolVar(&outputDebug, "v", outputDebug, "display debug output to ~/.statstee.log")
	flag.Parse()

	log.SetOutput(ioutil.Discard)
}

func main() {
	if outputDebug {
		logWriter, _ = os.OpenFile("statstee.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		log.SetFlags(log.Lmicroseconds | log.LstdFlags | log.Lshortfile)
		log.SetOutput(logWriter)
		defer logWriter.Close()
	}

	c := make(chan datagram.Metric, windowSize)
	router := router.New(c)

	go captureMetrics(router)
	go sniffStream(c)

	logIfError(views.Loop(router))
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
