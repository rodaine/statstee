package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/juju/ratelimit"
	"github.com/rodaine/statstee/datagram"
)

var (
	num       = 1
	rps int64 = 1000
)

func init() {
	flag.IntVar(&num, "n", num, "number of different metrics to generate")
	flag.Int64Var(&rps, "r", rps, "the metrics per second to send")
	flag.Parse()
}

func main() {
	limiter := ratelimit.NewBucketWithRate(float64(rps), rps)

	for n := runtime.NumCPU(); n >= 0; n-- {
		go beAJerk(limiter)
	}

	waitForSignal()
}

func beAJerk(limiter *ratelimit.Bucket) {
	sender, _ := datagram.NewSender("localhost", 8125)
	for {
		limiter.Wait(1)
		sender.Send(datagram.Metric{
			Type:       datagram.Histogram,
			Name:       "jerks." + strconv.Itoa(rand.Intn(num)),
			Value:      rand.Float64(),
			SampleRate: 1,
		})
	}
}

func waitForSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	log.Println("kill signal received")
}
