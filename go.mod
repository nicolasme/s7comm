module github.com/nicolasme/s7comm

go 1.16

require (
	github.com/influxdata/telegraf v1.18.3
	github.com/robinson/gos7 v0.0.0-20201216125248-2dd72fe148d3
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
)

replace (
	github.com/nicolasme/s7comm/plugins/inputs/s7comm => ./plugings/inputs/s7comm
)
