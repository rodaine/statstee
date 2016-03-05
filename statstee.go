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
	"github.com/rodaine/statstee/streams"
	"github.com/rodaine/statstee/views"
	"golang.org/x/net/context"
)

const (
	logFile = "statstee.log"
)

var (
	logWriter   io.WriteCloser
	streamError error

	deviceInterface string = streams.LoopbackAbbr
	sniffedPort     int    = streams.DefaultStatsDPort
	outputDebug     bool   = false
	listenMode      bool   = false
	captureMode     bool   = false
)

func init() {
	flag.StringVar(&deviceInterface, "d", deviceInterface, "network device to capture on")
	flag.IntVar(&sniffedPort, "p", sniffedPort, "port to capture on")
	flag.BoolVar(&listenMode, "l", listenMode, "force listen mode, error if the port cannot be bound")
	flag.BoolVar(&captureMode, "c", captureMode, "force capture mode, even if StatsD is not present")
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

	mode := streams.DefaultMode
	switch {
	case listenMode:
		mode = streams.ListenMode
	case captureMode:
		mode = streams.CaptureMode
	}
	stream, err := streams.ResolveStream(mode, deviceInterface, sniffedPort)
	fatalIfError(err)

	parser := datagram.NewParser()
	router := router.New(parser.Chan())

	go captureMetrics(router)
	go parser.Parse(stream.Chan())
	go stream.Listen(context.TODO())

	fatalIfError(views.Loop(router))
	fatalIfError(streamError)
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
