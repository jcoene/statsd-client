package statsd

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jcoene/gologger"
)

var NL = []byte("\n")

type Client struct {
	Addr   string
	Prefix string
	Proto  string
	Delay  time.Duration
	Debug  bool

	buf []byte
	mu  sync.Mutex
	log *logger.Logger
}

func NewClient(addr string, prefix string) (c *Client) {
	c = &Client{
		Addr:   addr,
		Prefix: prefix,
		Proto:  "udp",
		Delay:  1 * time.Second,
		log:    logger.NewDefaultLogger("statsd"),
	}

	go c.poll()

	return
}

func (c *Client) SetTCP() {
	c.Proto = "tcp"
}

func (c *Client) SetDebug(b bool) {
	c.Debug = b
}

func (c *Client) SetDelay(s string) {
	dur, err := time.ParseDuration(s)
	if err != nil {
		c.log.Error("invalid duration '%s': %s", s, err)
		return
	}

	c.Delay = dur
}

func (c *Client) Count(k string, d int64) {
	m := fmt.Sprintf("%s:%d|c", c.prefix(k), d)
	c.send([]byte(m))
}

func (c *Client) Inc(k string, d int64) {
	c.Count(k, d)
}

func (c *Client) Dec(k string, d int64) {
	c.Count(k, -d)
}

func (c *Client) Gauge(k string, v float64) {
	m := fmt.Sprintf("%s:%.3f|g", c.prefix(k), v)
	c.send([]byte(m))
}

func (c *Client) Measure(k string, v float64) {
	m := fmt.Sprintf("%s:%.3f|ms", c.prefix(k), v)
	c.send([]byte(m))
}

func (c *Client) Timing(k string, v float64) {
	c.Measure(k, v)
}

func (c *Client) MeasureDur(k string, dur time.Duration) {
	v := float64(dur) / float64(time.Millisecond)
	m := fmt.Sprintf("%s:%.3f|ms", c.prefix(k), v)
	c.send([]byte(m))
}

func (c *Client) send(data []byte) {
	if c.Debug {
		fmt.Println("StatsD:", string(data))
		return
	}

	c.mu.Lock()
	c.buf = append(c.buf, data...)
	c.buf = append(c.buf, NL...)
	defer c.mu.Unlock()

	return
}

func (c *Client) prefix(k string) string {
	if c.Prefix == "" {
		return k
	}

	return fmt.Sprintf("%s.%s", c.Prefix, k)
}

func (c *Client) poll() {
	var buf []byte

	for {
		time.Sleep(c.Delay)

		// Snapshot the buffer
		c.mu.Lock()
		buf = c.buf
		c.buf = []byte{}
		c.mu.Unlock()

		if len(buf) == 0 {
			continue
		}

		if err := c.push(buf); err != nil {
			if c.Proto == "tcp" {
				c.log.Error("unable to push %s %s: %s", c.Proto, c.Addr, err)

				// Replace the failed send buffer
				c.mu.Lock()
				c.buf = append(buf, c.buf...)
				c.mu.Unlock()
			}
		}
	}
}

func (c *Client) push(buf []byte) error {
	conn, err := net.Dial(c.Proto, c.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Split the buffer into multiple packets for UDP
	var packets [][]byte
	switch c.Proto {
	case "tcp":
		packets = [][]byte{buf}
	case "udp":
		packets = bytes.Split(buf, NL)
	}

	// Write the packets
	for _, p := range packets {
		if len(p) == 0 {
			continue
		}

		n, err := conn.Write(p)
		if err != nil {
			return err
		}

		if n != len(p) {
			return fmt.Errorf("wrote %d of %d bytes", n, len(p))
		}
	}

	return nil
}
