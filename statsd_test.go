package statsd

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ListenOnce() (ch chan string) {
	ch = make(chan string)

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8005")
	if err != nil {
		ch <- "addr error"
		return
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		ch <- "socket error"
		return
	}

	go func(sock *net.UDPConn, ch chan string) {
		defer sock.Close()

		buf := make([]byte, 512)
		n, _, err := sock.ReadFrom(buf)
		if err != nil {
			ch <- fmt.Sprintf("read error: %s", err)
			return
		}

		ch <- string(buf[0:n])
	}(sock, ch)

	time.Sleep(100 * time.Millisecond)

	return
}

func CompareFrom(ch chan string, expect string, t *testing.T) {
	select {
	case res := <-ch:
		if res != expect {
			t.Errorf("expected %s, got %s", expect, res)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("timeout")
	}
}

func TestDebug(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")

	defer func() {
		// Clean up the listening socket
		cli.SetDebug(false)
		cli.Count("customers.new", 3)
		<-ch
	}()

	cli.SetDebug(true)
	defer cli.SetDebug(false)
	cli.Count("customers.new", 3)

	time.Sleep(500 * time.Millisecond)

	select {
	case s := <-ch:
		t.Errorf("expected nothing to be received in debug mode, got %s", s)
	default:
	}
}

func TestCount(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Count("customers.new", 3)
	CompareFrom(ch, "myapp.customers.new:3|c", t)
}

func TestInc(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Inc("invoices.received", 3<<30)
	CompareFrom(ch, "myapp.invoices.received:3221225472|c", t)
}

func TestDec(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Dec("customers.maintained", 60)
	CompareFrom(ch, "myapp.customers.maintained:-60|c", t)
}

func TestGauge(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Gauge("queue.default.depth", 342)
	CompareFrom(ch, "myapp.queue.default.depth:342|g", t)
}

func TestIncGauge(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.IncGauge("queue.default.depth", 5)
	CompareFrom(ch, "myapp.queue.default.depth:5|ig", t)
}

func TestMeasure(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Measure("web.response.duration", 142)
	CompareFrom(ch, "myapp.web.response.duration:142|ms", t)
}

func TestMeasureDurMs(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.MeasureDur("web.response.duration", 32*time.Millisecond)
	CompareFrom(ch, "myapp.web.response.duration:32|ms", t)
}

func TestMeasureDurSec(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.MeasureDur("job.hardjob.duration", 11*time.Second)
	CompareFrom(ch, "myapp.job.hardjob.duration:11000|ms", t)
}

func TestMeasureDurMin(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.MeasureDur("job.hardjob.duration", 11*time.Minute)
	CompareFrom(ch, "myapp.job.hardjob.duration:660000|ms", t)
}

func TestTiming(t *testing.T) {
	ch := ListenOnce()
	cli, _ := NewClient("127.0.0.1:8005", "myapp")
	cli.Timing("web.response.duration", 142)
	CompareFrom(ch, "myapp.web.response.duration:142|ms", t)
}
