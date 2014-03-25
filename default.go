package statsd

import (
	"errors"
	"time"
)

var defaultClient *Client

var ErrNoDefaultClientConfigured = errors.New("statsd: no default client configured")

func NewDefaultClient(addr string, prefix string) (err error) {
	defaultClient, err = NewClient(addr, prefix)
	return
}

func SetDebug(b bool) {
	defaultClient.SetDebug(b)
}

func Count(k string, d int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Count(k, d)
}

func Inc(k string, d int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Inc(k, d)
}

func Dec(k string, d int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Dec(k, d)
}

func Gauge(k string, v int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Gauge(k, v)
}

func Measure(k string, v int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Measure(k, v)
}

func Timing(k string, v int64) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.Timing(k, v)
}

func MeasureDur(k string, dur time.Duration) error {
	if defaultClient == nil {
		return ErrNoDefaultClientConfigured
	}

	return defaultClient.MeasureDur(k, dur)
}
