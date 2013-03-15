package gocollectd

import (
	"time"
	"bytes"
	"encoding/binary"
)

const (
	TypeCounter  = 0
	TypeGuage    = 1
	TypeDerive   = 2
	TypeAbsolute = 3
)

type Value interface {
	CollectdType() uint8
}

type Counter uint64

func (v Counter) CollectdType() uint8 {
	return TypeCounter
}

type Guage float64

func (v Guage) CollectdType() uint8 {
	return TypeGuage
}

type Derive int64

func (v Derive) CollectdType() uint8 {
	return TypeDerive
}

type Absolute uint64

func (v Absolute) CollectdType() uint8 {
	return TypeAbsolute
}

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

func (p Packet) Values() ([]Value, error) {
	r := make([]Value, len(p.DataTypes))

	var err error
	buf := bytes.NewBuffer(p.Bytes)
	for i, t := range p.DataTypes {
		switch t {
		case TypeCounter:
			var v Counter
			err = binary.Read(buf, binary.BigEndian, &v)
			r[i] = v
		case TypeGuage:
			var v Guage
			err = binary.Read(buf, binary.BigEndian, &v)
			r[i] = v
		case TypeDerive:
			var v Derive
			err = binary.Read(buf, binary.BigEndian, &v)
			r[i] = v
		case TypeAbsolute:
			var v Absolute
			err = binary.Read(buf, binary.BigEndian, &v)
			r[i] = v
		}
	}
	return r, err
}
