# StatsD Client

[![Build Status](https://travis-ci.org/jcoene/statsd-client.png)](https://travis-ci.org/jcoene/statsd-client)

This is a no-nonsense StatsD client written in Go. It reconnects automatically if the write socket fails, making it at least 1000% better than any other StatsD client available.

It doesn't support sampling, I'd welcome a good pull request.

Pair it with [statsd-librato server](https://github.com/jcoene/statsd-librato) for a simple and efficient pipeline to Librato.

## Usage

```go
import "github.com/jcoene/statsd-client"

client, err := statsd.NewClient("127.0.0.1:8125", "myapp")

client.Count("http.response.200.count", 1)         // myapp.http.response.200.count:1|c
client.Inc("http.response.200.count", 1)           // myapp.http.response.200.count:1|c
client.Dec("http.response.200.count", 1)           // myapp.http.response.200.count:-1|c
client.Measure("http.response.200.runtime", 50)    // myapp.http.response.200.runtime:50|ms
client.MeasureDur("job.runtime", 50 * time.Second) // myapp.job.runtime:50000|ms
client.Gauge("server.load", 5)                     // myapp.server.load:5|g
```

You can also establish a default client that can be shared between parts of your app:

```go
import "github.com/jcoene/statsd-client"

err := statsd.NewDefaultClient("127.0.0.1:8125", "myapp")

statsd.Count("http.response.200.count", 1)         // myapp.http.response.200.count:1|c
```

## License

MIT License, see LICENSE file.
