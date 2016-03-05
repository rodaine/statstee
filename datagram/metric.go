package datagram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MetricType string

const (
	Counter   MetricType = "c"
	Gauge     MetricType = "g"
	Histogram MetricType = "h"
	Timer     MetricType = "ms"
	Set       MetricType = "s"
	Unknown   MetricType = "?"

	SampleRatePrefix = byte('@')
	TagsPrefix       = byte('#')
)

var (
	MalformedMetricError = errors.New("the metric is malformed")

	DummyMetric = Metric{
		Name: "statstee",
		Type: Unknown,
	}

	Prefix = map[MetricType]string{
		Histogram: "H",
		Timer:     "T",
		Counter:   "C",
		Gauge:     "G",
		Set:       "S",
		Unknown:   "?",
	}
)

type Metric struct {
	Name       string
	Value      float64
	Type       MetricType
	SampleRate float64
	Tags       []string
}

func ParseMetric(raw string) (m Metric, err error) {
	parts := strings.Split(raw, "|")
	if len(parts) < 2 {
		return m, MalformedMetricError
	}

	kv := strings.Split(parts[0], ":")
	if len(kv) != 2 {
		return m, MalformedMetricError
	}

	m.Name = kv[0]
	if m.Value, err = strconv.ParseFloat(kv[1], 64); err != nil {
		return
	}

	if parts[1] == "" {
		return m, MalformedMetricError
	}
	m.Type = MetricType(parts[1])

	m.SampleRate = 1

	if len(parts) > 2 {
		for _, part := range parts[2:] {
			switch part[0] {
			case SampleRatePrefix:
				if m.SampleRate, err = strconv.ParseFloat(part[1:], 64); err != nil {
					return
				}
			case TagsPrefix:
				m.Tags = strings.Split(part[1:], ",")
			default:
				return m, MalformedMetricError
			}
		}
	}

	return
}

func (m Metric) String() string {
	out := fmt.Sprintf("%s:%g|%s", m.Name, m.Value, m.Type)

	if m.SampleRate != 1 {
		out += fmt.Sprintf("|@%g", m.SampleRate)
	}

	if len(m.Tags) > 0 {
		out += fmt.Sprintf("|#%s", strings.Join(m.Tags, ","))
	}

	return out
}

func (m Metric) TypePrefix() string {
	if c, ok := Prefix[m.Type]; ok {
		return c
	}

	return Prefix[Unknown]
}
