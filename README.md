# statstee

_The Cross-Platform Proxyless Command Line UI for StatsD Metrics_

![Demo Animation of UI](/example/images/demo.gif?raw=true)

The `statstee` utility collects and plots metrics emitted from other processes on a host without interrupting or proxying their flow to the [StatsD daemon][statsd] or subsequent backend services.

`statstee` relies on [libpcap][libpcap] ([WinPcap][winpcap] for Windows) to capture the UDP metric datagrams off the network. Values are bucketed in one-second intervals and plots are shown depending on the metric datatype. Current and moving averages are also provided for each graph.

## Features

- [x] Cross-Platform (OSX / Linux / Windows)
- [x] Proxyless capture mode (daemon and processes can still use default host/port)
- [x] Listener mode if no daemon is present (prevents packet loss)
- [x] [Configurable](#usage) with sensible defaults (loopback interface, port 8125)
- [x] Supports all StatsD metric data-types ([gauge](#gauge), [counter](#counter), [set](#set), [timing](#timing--histogram-datadog))
- [x] Supports DataDog histogram metrics (treated like standard StatsD timing)
- [x] Graphed time-series for each metric type
- [x] Current value with 1-, 5- & 10-minute moving averages (EWMA)

### Planned Features

- [ ] Filter/search support
- [ ] Sample rate support
- [ ] DataDog [events][dd-events] and [service checks][dd-checks]
- [ ] Save and/or replay metric data

### Gauge

* Current Value

![Gauge Metrics](/example/images/gauge.png?raw=true)

### Counter

* Count / RPS
* Cumulative Count

![Counter Metrics](/example/images/count.png?raw=true)

### Set

* Unique Count / Unique RPS
* Percent Unique

![Set Metrics](/example/images/set.png?raw=true)

### Timing / Histogram (DataDog)

* Count / RPS
* Median
* 75<sup>th</sup> Percentile
* 95<sup>th</sup> Percentile

![Timing Metrics](/example/images/timing.png?raw=true)

## Install

Static binaries for many major platforms are forthcoming. For now, please refer to the [development section](#development) for installation instructions.

## Usage

```
→ statstee -h
Usage of dist/statstee:
  -c bool
     force capture mode, even if StatsD is not present
  -d string
      network device to capture on (default "_first_loopback_")
  -l bool
      force listen mode, error if the port cannot be bound
  -p int
      port to capture on (default 8125)
  -v bool
      display debug output to statstee.log
```

**NB:** You will likely need to run `statstee` as root or with `sudo` in order to snoop on the network traffic.

## Development

_Requires Go 1.5 or Higher_

#### Install libpcap for your OS & architecture

* **OSX:** `brew install libpcap`
* **Ubuntu/Debian:** `apt-get install libpcap-dev`
* **Windows:** [run installer][winpcap-install]

#### Download statstee

`go get github.com/rodaine/statstee`

#### Run tests

`script/bootstrap && script/test`

#### Build executables

`script/bootstrap && script/build`

#### Run executables

* `dist/statstee -v` - This app, with logging to `./statstee.log`!
* `dist/statter` - Demo app for experimenting (outputs all metric types)
* `dist/jerks -n 1 -r 1000` - Load testing app (`n` metrics at `r` rps)

## License

The MIT License (MIT)

Copyright (c) 2016 Chris Roche

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[statsd]:          https://github.com/etsy/statsd
[libpcap]:         http://www.tcpdump.org/
[winpcap]:         http://www.winpcap.org/
[winpcap-install]: http://www.winpcap.org/install/default.htm
[dd]:              https://www.datadoghq.com/
[dd-events]:       http://docs.datadoghq.com/guides/dogstatsd/#events-1
[dd-checks]:       http://docs.datadoghq.com/guides/dogstatsd/#service-checks
