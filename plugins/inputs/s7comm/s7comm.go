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

	Nodes []NodeSettings `toml:"nodes"`
	// Nodes []string `toml:"nodes"`

	Log telegraf.Logger `toml:"-"`

	// internal values
	handler *gos7.TCPClientHandler
	client  gos7.Client
	helper  gos7.Helper
}

type NodeSettings struct {
	Name string `toml:"name"`
	Type string `toml:"type"`
}

// type NodeSettings struct {
// 	Name   string `toml:"name"`
// 	Area   int    `toml:"area"`
// 	Db     int    `toml:"db"`
// 	Start  int    `toml:"start"`
// 	Amount int    `toml:"amount"`
// 	Length int    `toml:"length"`
// }

func (s *S7Comm) Connect() error {
	s.handler = gos7.NewTCPClientHandler(s.Endpoint, s.Rack, s.Slot)
	s.handler.Timeout = time.Duration(s.Timeout)
	s.handler.IdleTimeout = time.Duration(s.IdleTimeout)

	err := s.handler.Connect()
	if err != nil {
		return err
	}

	defer s.handler.Close()

	s.client = gos7.NewClient(s.handler)

	s.Log.Info("Connexion successfull")

	s.InitVars()

	return nil
}

func (s *S7Comm) Read() {

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
		# test1 = 1
		# nodes = [
		#	{name="test", area=131, db=1, start=0, amount=1, length=2},
		#	{name="test2", area=131, db=1, start=2, amount=1, length=2},
		# ]
`
}

type S7Result struct {
	Buf [4]byte
	Err string
}

func (s *S7Comm) InitVars() {

	buf := make([]byte, 4)

	for _, node := range s.Nodes {
		res, err := s.client.Read(node.Name, buf)

		if err != nil {
			s.Log.Error(err)
		} else {
			if node.Type != "real" {
				s.Log.Infof("Read %s : %d", node.Name, res)
			} else {
				t := s.helper.GetRealAt(buf, 0)
				s.Log.Infof("Read %s : %.2f", node.Name, t)
			}
		}
	}
}

// func (s *S7Comm) InitVars() {

// 	var items = []gos7.S7DataItem{}
// 	var results = []S7Result{}

// 	s.nodes = append(s.nodes, NodeSettings{
// 		Name:   "test",
// 		Area:   0x84,
// 		Db:     1,
// 		Start:  0,
// 		Amount: 1,
// 		Length: 4,
// 	})

// 	s.nodes = append(s.nodes, NodeSettings{
// 		Name:   "test2",
// 		Area:   0x84,
// 		Db:     1,
// 		Start:  2,
// 		Amount: 1,
// 		Length: 4,
// 	})

// 	for i, node := range s.nodes {
// 		s.Log.Infof("Read name : %s", node.Name)
// 		results = append(results, S7Result{Buf: [4]byte{}, Err: ""})
// 		items = append(items, gos7.S7DataItem{
// 			Area:     node.Area,
// 			WordLen:  node.Length,
// 			DBNumber: node.Db,
// 			Start:    node.Start,
// 			Amount:   node.Amount,
// 			Data:     results[i].Buf[:],
// 			Error:    results[i].Err,
// 		})
// 	}

// 	err := s.client.AGReadMulti(items, 2)
// 	if err != nil {
// 		s.Log.Error(err)
// 	}

// 	for _, res := range results {
// 		s.Log.Infof("Read buf : %d", res.Buf)

// 		t := binary.BigEndian.Uint16(res.Buf[:])
// 		s.Log.Infof("Read int : %d", t)

// 		// t2 := binary.
// 	}
// }

func (s *S7Comm) Gather(a telegraf.Accumulator) error {
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
