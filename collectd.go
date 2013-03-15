// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

// gocollectd parses the collectd binary protocol.
package gocollectd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	TypeCounter  = 0
	TypeGuage    = 1
	TypeDerive   = 2
	TypeAbsolute = 3
)

// A collectd value is sent as a int64, uint64, or float64. Number provides
// some useful functions to help handle this flexibility.
type Number interface {
	// CollectdType gives the Collectd Type of this number.
	CollectdType() uint8

	// Float64 converts this number to a float to avoid type assertions.
	Float64()      float64
}

// A collectd Counter value
type Counter uint64

// CollectdType returns TypeCounter
func (v Counter) CollectdType() uint8 {
	return TypeCounter
}

// Float64 converts this number to a float to avoid type assertions.
func (v Counter) Float64() float64 {
	return float64(v)
}

// A collectd Guage value
type Guage float64

// CollectdType returns TypeGuage
func (v Guage) CollectdType() uint8 {
	return TypeGuage
}

// Float64 converts this number to a float to avoid type assertions.
func (v Guage) Float64() float64 {
	return float64(v)
}

// A collectd Derive value
type Derive int64

// CollectdType returns TypeDerive
func (v Derive) CollectdType() uint8 {
	return TypeDerive
}

// Float64 converts this number to a float to avoid type assertions.
func (v Derive) Float64() float64 {
	return float64(v)
}

// A collectd Absolute value
type Absolute uint64

// CollectdType returns TypeAbsolute
func (v Absolute) CollectdType() uint8 {
	return TypeAbsolute
}

// Float64 converts this number to a float to avoid type assertions.
func (v Absolute) Float64() float64 {
	return float64(v)
}

// A packet is a set of collectd values that were sent at once by a collectd
// plugin.
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

// TimeUnixNano returns the measurement time in nanoseconds since unix epoch.
func (p Packet) TimeUnixNano() int64 {
	// 1.0737... is 2^30 (collectds' subsecond interval) / 10^-9 (nanoseconds)
	return int64(float64(p.CdTime) / 1.073741824)
}

// TimeUnix returns the measurement time in seconds since unix epoch.
func (p Packet) TimeUnix() int64 {
	return int64(p.CdTime >> 30)
}

// Time returns the measurement time as a go time.
func (p Packet) Time() time.Time {
	return time.Unix(0, p.TimeUnixNano())
}

// ValueCount returns the number of values in this packet.
func (p Packet) ValueCount() int {
	return len(p.DataTypes)
}

// ValueBytes returns the raw bytes for each value.
func (p Packet) ValueBytes() [][]byte {
	b := make([][]byte, len(p.DataTypes))
	for i := range b {
		b[i] = p.Bytes[i*8 : 8+(i*8)]
	}
	return b
}

// ValueNumbers returns the values as Numbers
func (p Packet) ValueNumbers() ([]Number, error) {
	r := make([]Number, len(p.DataTypes))

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

// ValueNames attempts to reformat collectd's plugin/type/instance heirarchy
// into strings.
func (p Packet) ValueNames() []string {
	r := make([]string, len(p.DataTypes))
	for i := range p.DataTypes {
		var name string
		// todo: think of ways to make this not a compiled in hack
		// todo: collectd 4 uses different patterns for some plugins
		// https://collectd.org/wiki/index.php/V4_to_v5_migration_guide
		switch p.Plugin {
		case "df":
			name = fmt.Sprintf("df_%s_%s", p.PluginInstance, p.TypeInstance)
		case "interface":
			switch i {
			case 0:
				name = fmt.Sprintf("%s_%s_tx", p.Type, p.PluginInstance)
			case 1:
				name = fmt.Sprintf("%s_%s_rx", p.Type, p.PluginInstance)
			}
		case "load":
			switch i {
			case 0:
				name = "load1"
			case 1:
				name = "load5"
			case 2:
				name = "load15"
			}
		case "memory":
			name = fmt.Sprintf("memory_%s", p.TypeInstance)
		default:
			name = fmt.Sprintf("%s_%s_%s_%s_%d", p.Plugin, p.PluginInstance, p.Type, p.TypeInstance, i)
		}
		r[i] = name
	}
	return r
}
