package s7comm

import (
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/robinson/gos7"
)

// S7Comm
type S7Comm struct {
	MetricName string `toml:"name"`
	Endpoint   string `toml:"plc_ip"`
	Rack       int    `toml:"plc_rack"`
	Slot       int    `toml:"plc_slot"`

	Timeout     config.Duration `toml:"connect_timeout"`
	IdleTimeout config.Duration `toml:"request_timeout"`

	// RootNodes []NodeSettings  `toml:"nodes"`
	// Groups    []GroupSettings `toml:"group"`

	// // internal values
	handler *gos7.TCPClientHandler
	client  *gos7.Client
	// req    *ua.ReadRequest
	// opts   []opcua.Option
}

func (s7Client *S7Comm) Connect() error {
	s7Client.handler = gos7.NewTCPClientHandler(s7Client.Endpoint, s7Client.Rack, s7Client.Slot)
	s7Client.handler.Timeout = time.Duration(s7Client.Timeout)
	s7Client.handler.IdleTimeout = time.Duration(s7Client.IdleTimeout)

	err := s7Client.handler.Connect()
	if err != nil {
		return err
	}

	defer s7Client.handler.Close()

	return nil
}

func (s7Client *S7Comm) Read() {

}

func (s7Client *S7Comm) Stop() error {
	err := s7Client.handler.Close()

	return err
}

func (s7Client *S7Comm) Init() error {
	err := s7Client.Connect()
	return err
}

func (s7Client *S7Comm) SampleConfig() string {
	return `
  	## Generates random numbers
		[[inputs.s7comm]]
		# name = "S7300"
		# plc_ip = "192.168.10.58"
		# plc_rack = 0
		# plc_slot = 2
		# connect_timeout = 10s
		# request_timeout = 2s
`
}

func (s7Client *S7Comm) Gather(a telegraf.Accumulator) error {
	return nil
}

func (s7Client *S7Comm) Description() string {
	return "Read data from SIemens PLC using S7Go"
}

// Add this plugin to telegraf
func init() {
	inputs.Add("s7comm", func() telegraf.Input {
		return &S7Comm{
			MetricName:  "s7comm",
			Endpoint:    "192.168.10.58",
			Rack:        0,
			Slot:        2,
			Timeout:     config.Duration(5 * time.Second),
			IdleTimeout: config.Duration(10 * time.Second),
		}
	})
}
