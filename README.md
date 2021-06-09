# S7comm

S7comm is a [telegraf](https://github.com/influxdata/telegraf) external input plugin to gather data from Siemens PLC using the [gos7](https://github.com/robinson/gos7) library.

## Installation

Download the repo

```bash
git clone git@github.com:nicolasme/s7comm.git
```

Change the poll interval if needed in the cmd/main.go file

```golang
var pollInterval = flag.Duration("poll_interval", 10*time.Second, "how often to send metrics")
```

Build the binary

```bash
go build -o s7comm cmd/main.go
```

Create your plugin.config file

```toml
[[inputs.s7comm]]
	name = "S7300"
	plc_ip = "192.168.10.57"
	plc_rack = 0
	plc_slot = 1
    	connect_timeout = "10s"
    	request_timeout = "2s"
    	nodes = [{name= "test_int", address= "DB1.DBW0", type = "int"},
        	{name= "test_real", address= "DB1.DBD2",type = "real"},
        	{name= "test_bool", address= "DB1.DBX10.0",type = "bool"},
        	{name= "test_dint", address= "DB1.DBD12",type = "dint"},
        	{name= "test_uint", address= "DB1.DBW16",type = "uint"},
        	{name= "test_udint", address= "DB1.DBD18",type = "udint"}]
```

From here, you can already test the plugin with your config file.

```bash
s7comm -config plugin.conf
```

If everything is ok, you should see something like this

```bash
2021/06/09 10:37:27 I! Connexion successfull
test_int value=8056i 1623227848846884706
test_real value=403.14764404296875 1623227849849296642
```

## Telegraf configuration

To use the plugin with telegraf, add this configuration to your main telegraf.conf file. S7comm is an external plugin using [shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim/README.md) and [execd](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd). Go see their doc for more information.

```toml
[[inputs.execd]]
  command = ["/path/to/s7comm", "-config", "/path/to/s7comm.plugin.conf"]
  signal = "none"
  restart_delay = "10s"
  data_format = "influx"
```

## PLC configuration

S7-300 and S7-400 usually use rack 0 and slot 2 and dont require additional configuration.

S7-1200 and S7-1500 usually use rack 0 and slot 1 and you need to enable the PUT/GET operations in the hardware configuration of your PLC.

Be aware of security issue. Once S7 Communication is enabled in a CPU, there is no way to block communication with a partner device. This means that any device on the same network can read and write data to the CPU using the S7 Communication protocol. For this reason, I would recommend using the native OPC.UA server for the newer S7-1200 and S7-1500 PLCs. See the [OPC.UA](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/opcua) telegraf plugin.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
