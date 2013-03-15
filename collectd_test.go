package gocollectd

import (
	"reflect"
	"testing"
	"time"
)

var testPacket = Packet{"laptop.lan", "interface", "lo0", "if_octets", "", 1463827927039889790, []uint8{TypeDerive, TypeGuage, TypeDerive}, h2b("00 00 00 00 00 88 07 8b 41 cf 43 00 00 00 00 00 00 00 00 00 00 88 07 8c")}
var testDate = time.Date(2013, time.March, 14, 21, 19, 53, 804828672, time.UTC)

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

func TestPacketValues(t *testing.T) {
	result, err := testPacket.Values()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := []Value{
		Derive(8914827),
		Guage(1048969216),
		Derive(8914828),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v", expected, result)
	}
}
