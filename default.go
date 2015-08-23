package statsd

import (
	"os"
	"time"
)

var DefaultClient *Client

func init() {
	Init()
}

func Init() {
	addr := os.Getenv("STATSD_HOST")
	if addr == "" {
		addr = "127.0.0.1:8125"
	}

	NewDefaultClient(addr, os.Getenv("STATSD_PREFIX"))

	if s := os.Getenv("STATSD_PROTO"); s != "" {
		DefaultClient.Proto = s
	}

	if s := os.Getenv("STATSD_DELAY"); s != "" {
		DefaultClient.SetDelay(s)
	}
}

func NewDefaultClient(addr string, prefix string) {
	DefaultClient = NewClient(addr, prefix)
}

func SetDebug(b bool) {
	DefaultClient.SetDebug(b)
}

func SetDelay(s string) {
	DefaultClient.SetDelay(s)
}

func Count(k string, d int64) {
	DefaultClient.Count(k, d)
}

func Inc(k string, d int64) {
	DefaultClient.Inc(k, d)
}

func Dec(k string, d int64) {
	DefaultClient.Dec(k, d)
}

func Gauge(k string, v float64) {
	DefaultClient.Gauge(k, v)
}

func Measure(k string, v float64) {
	DefaultClient.Measure(k, v)
}

func Timing(k string, v float64) {
	DefaultClient.Timing(k, v)
}

func MeasureDur(k string, dur time.Duration) {
	DefaultClient.MeasureDur(k, dur)
}
