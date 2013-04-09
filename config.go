package esmc

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	OnMode = Mode(iota)
	OffMode
	DarkMode
)

type Mode int

func (m Mode) String() string {
	switch m {
	case OnMode:
		return "on"
	case OffMode:
		return "off"
	case DarkMode:
		return "dark"
	}

	panic(fmt.Sprintf("invalid cluster mode: %d", m))
}

type Config struct {
	Endpoints []string
	options   url.Values
}

// Build a new config struct from a specification string. The spec format is:
//
//		es://host[:port][,host[:port]...][?options]
//
// Where valid options are:
//
//		name = string [required, no default]
//		mode = on, dark, off [default: off]
//    ping_timeout = duration string [default: 250ms]
//    ping_interval = duration string [default: 10s]
//
func NewConfig(spec string) (Config, error) {
	uri, err := url.Parse(spec)
	if err != nil {
		return Config{}, err
	}

	options := uri.Query()
	uri.RawQuery = "" // strip query params

	if options.Get("name") == "" {
		return Config{}, fmt.Errorf("cluster must have a name: %s", spec)
	}

	endpoints := []string{}
	for _, host := range strings.Split(uri.Host, ",") {
		endpoint := &url.URL{Scheme: "http", Host: host}
		endpoints = append(endpoints, endpoint.String())
	}

	return Config{endpoints, options}, nil
}

// Returns the name of the cluster.
func (c Config) Name() string {
	return c.options.Get("name")
}

// Returns the cluster mode (OnMode, OffMode, or DarkMode).
func (c Config) Mode() Mode {
	switch c.options.Get("mode") {
	case "on":
		return OnMode
	case "dark":
		return DarkMode
	}

	return OffMode
}

// Defines how long to wait for a node to response to a ping before marking it
// as down. Default 250ms.
func (c Config) PingTimeout() time.Duration {
	timeout, err := time.ParseDuration(c.options.Get("ping_timeout"))

	if err == nil {
		return timeout
	}

	return 250 * time.Millisecond
}

// Defines how often to ping the defined endpoints. Default 10s.
func (c Config) PingInterval() time.Duration {
	interval, err := time.ParseDuration(c.options.Get("ping_interval"))

	if err == nil {
		return interval
	}

	return 10 * time.Second
}
