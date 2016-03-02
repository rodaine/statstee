package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PagerDuty/godspeed"
)

func main() {
	gs, _ := godspeed.NewDefault()
	gs.SetNamespace("statter")

	go histogram(gs)

	go timer(gs)
	go count(gs)
	go gauge(gs)
	go set(gs)
	waitForSignal()
}

func histogram(gs *godspeed.Godspeed) {
	for range time.NewTicker(time.Millisecond * 3).C {
		gs.Histogram("histogram", math.Sin(float64(time.Now().Unix())/math.Pi), nil)
	}
}

func count(gs *godspeed.Godspeed) {
	for range time.NewTicker(time.Millisecond * 5).C {
		gs.Count("count", math.Abs(math.Sin(float64(time.Now().Unix())/math.Pi))/1000, nil)
	}
}

func gauge(gs *godspeed.Godspeed) {
	val := 0.0
	for range time.NewTicker(time.Millisecond * 7).C {
		val += math.Abs(rand.NormFloat64()) / 1000
		gs.Gauge("gauge", val, nil)
	}
}

func set(gs *godspeed.Godspeed) {
	for range time.NewTicker(time.Millisecond).C {
		gs.Set("set", float64(rand.Intn(time.Now().Second()+1)), nil)
	}
}

func timer(gs *godspeed.Godspeed) {
	for range time.NewTicker(time.Millisecond * 3).C {
		gs.Timing("timer", math.Cos(float64(time.Now().Unix())/math.Pi), nil)
	}
}

func waitForSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	log.Println("kill signal received")
}
