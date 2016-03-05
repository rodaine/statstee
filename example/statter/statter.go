package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rodaine/statstee/datagram"
)

func main() {
	go histogram()
	go timer()
	go count()
	go gauge()
	go set()

	waitForSignal()
}

func histogram() {
	sender, _ := datagram.NewSender("localhost", 8125)
	for range time.NewTicker(time.Millisecond * 3).C {
		sender.Send(datagram.Metric{
			Type:       datagram.Histogram,
			Name:       "statter.histogram",
			Value:      math.Sin(float64(time.Now().Unix())/math.Pi) + rand.Float64(),
			SampleRate: 1,
		})
	}
}

func count() {
	sender, _ := datagram.NewSender("localhost", 8125)
	for range time.NewTicker(time.Millisecond * 5).C {
		sender.Send(datagram.Metric{
			Type:       datagram.Counter,
			Name:       "statter.count",
			Value:      math.Abs(math.Sin(float64(time.Now().Unix())/math.Pi)) / 10,
			SampleRate: 1,
		})
	}
}

func gauge() {
	sender, _ := datagram.NewSender("localhost", 8125)
	for range time.NewTicker(time.Millisecond * 7).C {
		sender.Send(datagram.Metric{
			Type:       datagram.Gauge,
			Name:       "statter.gauge",
			Value:      math.Sin(float64(time.Now().Unix())/10*math.Pi) + math.Cos(float64(time.Now().Second())/100*math.Pi),
			SampleRate: 1,
		})
	}
}

func set() {
	sender, _ := datagram.NewSender("localhost", 8125)
	for range time.NewTicker(time.Millisecond).C {
		sender.Send(datagram.Metric{
			Type:       datagram.Set,
			Name:       "statter.set",
			Value:      float64(rand.Intn(time.Now().Second() + 1)),
			SampleRate: 1,
		})
	}
}

func timer() {
	sender, _ := datagram.NewSender("localhost", 8125)
	for range time.NewTicker(time.Millisecond * 11).C {
		x := float64(time.Now().Unix()) / math.Pi
		sender.Send(datagram.Metric{
			Type:       datagram.Timer,
			Name:       "statter.timer",
			Value:      1 + math.Min(math.Sin(x), math.Cos(x)),
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
