package gocollectd

import(
	"time"
)

const (
	Counter  = 0
	Guage    = 1
	Derive   = 2
	Absolute = 3
)

type Value struct {
	Hostname       string
	Plugin         string
	PluginInstance string
	Type           string
	TypeInstance   string
	Number         uint16
	CdTime         uint64
	DataType       uint8
	Bytes          []byte
}

func (v Value) TimeUnixNano() int64 {
	// 1.0737... is 2^30 (collectds' subsecond interval) / 10^-9 (nanoseconds)
	return int64(float64(v.CdTime) / 1.073741824)
}

func (v Value) TimeUnix() int64 {
	return int64(v.CdTime >> 30)
}

func (v Value) Time() time.Time {
	return time.Unix(0, v.TimeUnixNano())
}
