package datagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric_Parse(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		Raw      string
		Error    bool
		Expected Metric
	}{
		{"", true, Metric{}},
		{"foo.bar", true, Metric{}},
		{"foo.bar:123", true, Metric{}},
		{"foo.bar:123|", true, Metric{}},
		{"foo.bar|h", true, Metric{}},
		{"foo.bar:fizz|h", true, Metric{}},
		{"foo.bar:123|h|@fizz", true, Metric{}},
		{"foo.bar:456|h|$invalid", true, Metric{}},

		{"foo.bar:123|h", false, Metric{Name: "foo.bar", Value: 123, Type: Histogram, SampleRate: 1}},
		{"foo.bar:123|h|@0.5", false, Metric{Name: "foo.bar", Value: 123, Type: Histogram, SampleRate: 0.5}},
		{"foo.bar:123|h|#tag1:val,tag2", false, Metric{Name: "foo.bar", Value: 123, Type: Histogram, SampleRate: 1, Tags: []string{"tag1:val", "tag2"}}},
		{"foo.bar:123|h|@0.5|#tag1:val,tag2", false, Metric{Name: "foo.bar", Value: 123, Type: Histogram, SampleRate: 0.5, Tags: []string{"tag1:val", "tag2"}}},
	}

	for _, test := range tests {
		actual, err := ParseMetric(test.Raw)

		if test.Error {
			is.Error(err, "%+v", test)
			continue
		}

		is.NoError(err, "%+v", test)
		is.EqualValues(test.Expected, actual, "%+v", test)
	}
}

func TestMetric_String(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		Expected string
		Metric   Metric
	}{
		{"foo.bar:123|c", Metric{Name: "foo.bar", Value: 123, Type: Counter, SampleRate: 1}},
		{"foo.bar:123|c|@0.5", Metric{Name: "foo.bar", Value: 123, Type: Counter, SampleRate: 0.5}},
		{"foo.bar:123|c|#tag1:val,tag2", Metric{Name: "foo.bar", Value: 123, Type: Counter, SampleRate: 1, Tags: []string{"tag1:val", "tag2"}}},
		{"foo.bar:123|c|@0.5|#tag1:val,tag2", Metric{Name: "foo.bar", Value: 123, Type: Counter, SampleRate: 0.5, Tags: []string{"tag1:val", "tag2"}}},
	}

	for _, test := range tests {
		is.Equal(test.Expected, test.Metric.String(), "%#v", test.Metric)
	}
}

func TestMetric_Prefix(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	tests := []struct {
		Type   MetricType
		Prefix string
	}{
		{Histogram, "H"},
		{Timer, "T"},
		{Counter, "C"},
		{Gauge, "G"},
		{Set, "S"},
		{Unknown, "?"},
		{"foobar", "?"},
	}

	for _, test := range tests {
		m := Metric{Type: test.Type}
		is.Equal(test.Prefix, m.TypePrefix(), "%#v", test)
	}
}
