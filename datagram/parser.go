package datagram

import "log"

type Parser interface {
	Parse(data <-chan []byte)
	Chan() <-chan Metric
}

type parser struct {
	c chan Metric
}

func NewParser() Parser {
	return &parser{make(chan Metric, 1000)}
}

func (p *parser) Parse(data <-chan []byte) {
	defer close(p.c)
	for raw := range data {
		m, err := ParseMetric(string(raw))
		if err != nil {
			log.Printf("unable to parse datagram: %v", err)
			continue
		}
		p.c <- m
	}
}

func (p *parser) Chan() <-chan Metric {
	return p.c
}
