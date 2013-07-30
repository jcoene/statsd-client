package statsd

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	Addr   string
	Prefix string
	nc     net.Conn
	mu     sync.Mutex
	rw     *bufio.ReadWriter
}

func NewClient(addr string, prefix string) (c *Client, err error) {
	c = &Client{
		Addr:   addr,
		Prefix: prefix,
	}

	err = c.redial()

	return
}

func (c *Client) Count(k string, d int64) error {
	m := fmt.Sprintf("%s:%d|c", c.prefix(k), d)
	return c.send([]byte(m))
}

func (c *Client) Inc(k string, d int64) error {
	return c.Count(k, d)
}

func (c *Client) Dec(k string, d int64) error {
	return c.Count(k, -d)
}

func (c *Client) Gauge(k string, v int64) error {
	m := fmt.Sprintf("%s:%d|g", c.prefix(k), v)
	return c.send([]byte(m))
}

func (c *Client) Measure(k string, v int64) error {
	m := fmt.Sprintf("%s:%d|ms", c.prefix(k), v)
	return c.send([]byte(m))
}
func (c *Client) Timing(k string, v int64) error {
	return c.Measure(k, v)
}

func (c *Client) redial() (err error) {
	c.nc = nil
	if c.nc, err = net.Dial("udp", c.Addr); err != nil {
		return
	}

	c.rw = bufio.NewReadWriter(bufio.NewReader(c.nc), bufio.NewWriter(c.nc))

	return
}

func (c *Client) send(data []byte) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err = c.rw.Write(data)
	if err != nil {
		fmt.Println("redialing")
		c.redial()
	}

	err = c.rw.Flush()

	return
}

func (c *Client) prefix(k string) string {
	if c.Prefix == "" {
		return k
	}

	return fmt.Sprintf("%s.%s", c.Prefix, k)
}