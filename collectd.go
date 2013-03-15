package gocollectd

import (
	"time"
)

const (
	Counter  = 0
	Guage    = 1
	Derive   = 2
	Absolute = 3
)

type Packet struct {
	Hostname       string
	Plugin         string
	PluginInstance string
	Type           string
	TypeInstance   string
	CdTime         uint64
	DataTypes      []uint8
	Bytes          []byte
}

func (p Packet) TimeUnixNano() int64 {
	// 1.0737... is 2^30 (collectds' subsecond interval) / 10^-9 (nanoseconds)
	return int64(float64(p.CdTime) / 1.073741824)
}

func (p Packet) TimeUnix() int64 {
	return int64(p.CdTime >> 30)
}

func (p Packet) Time() time.Time {
	return time.Unix(0, p.TimeUnixNano())
}

func (p Packet) ValueBytes() [][]byte {
	b := make([][]byte, len(p.DataTypes))
	for i := range b {
		b[i] = p.Bytes[i*8 : 8+(i*8)]
	}
	return b
}
