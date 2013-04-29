// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

// gocollectd parses the collectd binary protocol.
package gocollectd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	TypeCounter  = 0
	TypeGauge    = 1
	TypeDerive   = 2
	TypeAbsolute = 3
)

// A collectd value is sent as a int64, uint64, or float64. Number provides
// some useful functions to help handle this flexibility.
type Number interface {
	// CollectdType gives the Collectd Type of this number.
	CollectdType() uint8

	// Float64 converts this number to a float to avoid type assertions.
	Float64() float64
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

// A collectd Gauge value
type Gauge float64

// CollectdType returns TypeGauge
func (v Gauge) CollectdType() uint8 {
	return TypeGauge
}

// Float64 converts this number to a float to avoid type assertions.
func (v Gauge) Float64() float64 {
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

// A Value is what a packet contains.
type Value struct {
	dataType  uint8
	bytes     []byte
}

// This value's raw bytes.
func (v Value) Bytes() []byte {
	return v.bytes
}

// Number converts this value into a Number.
func (v Value) Number() (Number, error) {
	r := bytes.NewReader(v.bytes)
	return byteReaderToNumber(v.dataType, r)
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
	CdInterval     uint64
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

// Values returns the data in this packet as a []Value.
func (p Packet) Values() ([]Value) {
	values := make([]Value, len(p.DataTypes))
	for i := range values {
		values[i] = Value{p.DataTypes[i], p.Bytes[i*8 : 8+(i*8)] }
	}
	return values
}

// ValueNumbers returns the values as Numbers
func (p Packet) ValueNumbers() ([]Number, error) {
	r := make([]Number, len(p.DataTypes))

	var err error
	reader := bytes.NewReader(p.Bytes)
	for i, t := range p.DataTypes {
		r[i], err = byteReaderToNumber(t, reader)
		if err != nil {
			return []Number{}, err
		}
	}
	return r, nil
}

// ValueNames attempts to reformat collectd's plugin/type/instance heirarchy
// into a string for this packet.
func (p Packet) Name() (name string) {
	// todo: think of ways to make this not a compiled in hack
	// todo: collectd 4 uses different patterns for some plugins
	// https://collectd.org/wiki/index.php/V4_to_v5_migration_guide
	switch p.Plugin {
	case "df":
		name = fmt.Sprintf("df_%s_%s", p.PluginInstance, p.TypeInstance)
	case "interface":
		name = fmt.Sprintf("%s_%s", p.Type, p.PluginInstance)
	case "load":
		name = "load"
	case "memory":
		name = fmt.Sprintf("memory_%s", p.TypeInstance)
	default:
		name = fmt.Sprintf("%s_%s_%s_%s", p.Plugin, p.PluginInstance, p.Type, p.TypeInstance)
	}
	return name
}

// ValueNames attempts to reformat collectd's plugin/type/instance heirarchy
// into a strings for each value in this packet.
func (p Packet) ValueNames() []string {
	r := make([]string, len(p.DataTypes))
	for i := range p.DataTypes {
		name := p.Name()
		var valueName string
		switch {
		case p.Plugin == "df" && i == 0:
			valueName = ""
		case p.Plugin == "memory" && i == 0:
			valueName = ""
		case p.Plugin == "interface" && i == 0:
			valueName = "tx"
		case p.Plugin == "interface" && i == 1:
			valueName = "rx"
		case p.Plugin == "load" && i == 0:
			valueName = "1"
		case p.Plugin == "load" && i == 1:
			valueName = "5"
		case p.Plugin == "load" && i == 2:
			valueName = "15"
		default:
			valueName = strconv.FormatInt(int64(i), 10)
		}
		if valueName == "" {
			r[i] = name
		} else {
			r[i] = fmt.Sprintf("%s_%s", name, valueName)
		}
	}
	return r
}

func byteReaderToNumber(dataType uint8, reader *bytes.Reader) (n Number, err error) {
	switch dataType {
	case TypeCounter:
		var v Counter
		err = binary.Read(reader, binary.BigEndian, &v)
		return v, err
	case TypeGauge:
		var v Gauge
		err = binary.Read(reader, binary.BigEndian, &v)
		return v, err
	case TypeDerive:
		var v Derive
		err = binary.Read(reader, binary.BigEndian, &v)
		return v, err
	case TypeAbsolute:
		var v Absolute
		err = binary.Read(reader, binary.BigEndian, &v)
		return v, err
	}
	return nil, errors.New("unknown value type")
}
