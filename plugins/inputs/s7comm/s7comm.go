package s7comm

import (
	"github.com/influxdata/telegraf/config"
)

// S7Comm
type S7Comm struct {
	MetricName string `toml:"name"`
	Endpoint   string `toml:"endpoint"`
	Rack       string `toml:"security_policy"`
	Slot       string `toml:"security_mode"`

	Timeout     config.Duration `toml:"connect_timeout"`
	IdleTimeout config.Duration `toml:"request_timeout"`

	// RootNodes []NodeSettings  `toml:"nodes"`
	// Groups    []GroupSettings `toml:"group"`

	// // internal values
	// client *TCPClientHandler
	// req    *ua.ReadRequest
	// opts   []opcua.Option
}

func init() {
}
