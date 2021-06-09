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

	Nodes []NodeSettings  `toml:"nodes"`
	Log   telegraf.Logger `toml:"-"`

	// internal values
	handler *gos7.TCPClientHandler
	client  gos7.Client
	helper  gos7.Helper
}

type NodeSettings struct {
	Name    string `toml:"name"`
	Address string `toml:"address"`
	Type    string `toml:"type"`
}

func (s *S7Comm) Connect() error {
	s.handler = gos7.NewTCPClientHandler(s.Endpoint, s.Rack, s.Slot)
	s.handler.Timeout = time.Duration(s.Timeout)
	s.handler.IdleTimeout = time.Duration(s.IdleTimeout)

	err := s.handler.Connect()
	if err != nil {
		return err
	}

	// defer s.handler.Close()

	s.client = gos7.NewClient(s.handler)

	s.Log.Info("Connexion successfull")

	return nil
}

func (s *S7Comm) Stop() error {
	err := s.handler.Close()

	return err
}

func (s *S7Comm) Init() error {
	err := s.Connect()
	return err
}

func (s *S7Comm) SampleConfig() string {
	return `
  	## Generates random numbers
		[[inputs.s7comm]]
		# name = "S7300"
		# plc_ip = "192.168.10.57"
		# plc_rack = 1
		# plc_slot = 2
		# connect_timeout = 10s
		# request_timeout = 2s
		# nodes = [{name= "DB1.DBW0", type = "int"}, 
        # {name= "DB1.DBD2", type = "real"},
        # {name= "DB1.DBD6", type = "real"}, 
        # {name= "DB1.DBX10.0", type = "bool"}, 
        # {name= "DB1.DBD12", type = "dint"}, 
        # {name= "DB1.DBW16", type = "uint"}, 
        # {name= "DB1.DBD18", type = "udint"}, 
        # {name= "DB1.DBD22", type = "time"}]
`
}

func (s *S7Comm) Gather(a telegraf.Accumulator) error {

	for _, node := range s.Nodes {

		var fields map[string]interface{}
		buf := make([]byte, 4)

		_, err := s.client.Read(node.Address, buf)

		if err != nil {
			s.Log.Error(err)
		} else {
			switch dataType := node.Type; dataType {
			case "bool":
				var res bool
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "byte":
				var res byte
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "word":
				var res uint16
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "dword":
				var res uint32
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "int":
				var res int16
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "dint":
				var res int32
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "uint":
				var res uint16
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "udint":
				var res uint32
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "real":
				var res float32
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			case "time":
				var res uint32
				s.helper.GetValueAt(buf, 0, &res)
				fields = map[string]interface{}{node.Name: res}
			default:
				s.Log.Error("Type error - ", node.Type)
				fields = nil
			}

			if fields != nil {
				a.AddFields(s.MetricName, fields, nil)
			}
		}
	}

	return nil
}

func (s *S7Comm) Description() string {
	return "Read data from Siemens PLC using S7 protocol with S7Go"
}

// Add this plugin to telegraf
func init() {
	inputs.Add("s7comm", func() telegraf.Input {
		return &S7Comm{
			MetricName:  "s7comm",
			Endpoint:    "192.168.10.57",
			Rack:        0,
			Slot:        1,
			Timeout:     config.Duration(5 * time.Second),
			IdleTimeout: config.Duration(10 * time.Second),
			Nodes:       nil,
		}
	})
}
