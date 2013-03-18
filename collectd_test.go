// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package gocollectd

import (
	"reflect"
	"testing"
	"time"
)

var testPacket = Packet{"laptop.lan", "fake", "", "", "", 1463827927039889790, []uint8{TypeDerive, TypeGuage, TypeDerive}, h2b("00 00 00 00 00 88 07 8b 41 cf 43 00 00 00 00 00 00 00 00 00 00 88 07 8c")}
var testDate = time.Date(2013, time.March, 14, 21, 19, 53, 804828672, time.UTC)
var testValue = Value{TypeGuage, h2b("41 cf 43 00 00 00 00 00")}

func TestValueBytes(t *testing.T) {
	result := testValue.Bytes()
	expected := testValue.bytes
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v", expected, result)
	}
}

func TestValueNumber(t *testing.T) {
	result, err := testValue.Number()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != Guage(1048969216) {
		t.Errorf("expected 1048969216 got %v", result)
	}
}

func TestPacketTime(t *testing.T) {
	result := testPacket.Time()
	if !result.Equal(testDate) {
		t.Errorf("expected %v, got %v", testDate, result)
	}
}

func TestPacketSeconds(t *testing.T) {
	result := testPacket.TimeUnix()
	expected := testDate.Unix()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPacketNanoSeconds(t *testing.T) {
	result := testPacket.TimeUnixNano()
	expected := testDate.UnixNano()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestPacketValueCount(t *testing.T) {
	result := testPacket.ValueCount()
	if result != 3 {
		t.Errorf("expected test packet to have 3 values, got %d", result)
	}
}

func TestPacketValueBytes(t *testing.T) {
	result := testPacket.ValueBytes()
	expected := [][]byte{
		h2b("00 00 00 00 00 88 07 8b"),
		h2b("41 cf 43 00 00 00 00 00"),
		h2b("00 00 00 00 00 88 07 8c"),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v", expected, result)
	}
}

func TestPacketValueNumbers(t *testing.T) {
	result, err := testPacket.ValueNumbers()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := []Number{
		Derive(8914827),
		Guage(1048969216),
		Derive(8914828),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v", expected, result)
	}
}

func TestPacketValues(t *testing.T) {
	result := testPacket.Values()
	expected := []Value{
		{TypeDerive, h2b("00 00 00 00 00 88 07 8b") },
		{TypeGuage,  h2b("41 cf 43 00 00 00 00 00") },
		{TypeDerive, h2b("00 00 00 00 00 88 07 8c") },
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v", expected, result)
	}
}

func TestPacketName(t *testing.T) {
	tests := []struct {
		packet Packet
		name   string
	}{
		{
			Packet{"laptop.lan", "interface", "lo0", "if_octets", "", 1463827927249453056, []uint8{TypeDerive, TypeDerive}, []byte{}},
			"if_octets_lo0",
		},
		{
			Packet{"laptop.lan", "memory", "", "memory", "wired", 1463827927249453056, []uint8{TypeGuage}, []byte{}},
			"memory_wired",
		},
		{
			Packet{"laptop.lan", "load", "", "load", "wired", 1463827927249453056, []uint8{TypeGuage, TypeGuage, TypeGuage}, []byte{}},
			"load",
		},
		{
			Packet{"laptop.lan", "df", "root", "df_complex", "used", 1463827927249453056, []uint8{TypeGuage}, []byte{}},
			"df_root_used",
		},
		{
			Packet{"laptop.lan", "plugin", "some", "thing", "here", 1463827927249453056, []uint8{TypeGuage, TypeGuage}, []byte{}},
			"plugin_some_thing_here",
		},
	}

	for i, tst := range tests {
		result := tst.packet.Name()
		if tst.name != result {
			t.Errorf("%i: expected\n%v\ngot\n%v", i, tst.name, result)
		}
	}

}

func TestPacketValueNames(t *testing.T) {
	tests := []struct {
		packet Packet
		names  []string
	}{
		{
			Packet{"laptop.lan", "interface", "lo0", "if_octets", "", 1463827927249453056, []uint8{TypeDerive, TypeDerive}, []byte{}},
			[]string{"if_octets_lo0_tx", "if_octets_lo0_rx"},
		},
		{
			Packet{"laptop.lan", "memory", "", "memory", "wired", 1463827927249453056, []uint8{TypeGuage}, []byte{}},
			[]string{"memory_wired"},
		},
		{
			Packet{"laptop.lan", "load", "", "load", "wired", 1463827927249453056, []uint8{TypeGuage, TypeGuage, TypeGuage}, []byte{}},
			[]string{"load_1", "load_5", "load_15"},
		},
		{
			Packet{"laptop.lan", "df", "root", "df_complex", "used", 1463827927249453056, []uint8{TypeGuage}, []byte{}},
			[]string{"df_root_used"},
		},
		{
			Packet{"laptop.lan", "plugin", "some", "thing", "here", 1463827927249453056, []uint8{TypeGuage, TypeGuage}, []byte{}},
			[]string{"plugin_some_thing_here_0", "plugin_some_thing_here_1"},
		},
	}

	for i, tst := range tests {
		result := tst.packet.ValueNames()
		if !reflect.DeepEqual(result, tst.names) {
			t.Errorf("%i: expected\n%v\ngot\n%v", i, tst.names, result)
		}
	}

}
