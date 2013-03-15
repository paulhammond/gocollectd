package gocollectd

import (
	"testing"
	"time"
)

var testValue = Value{"laptop.lan", "memory", "", "", "wired", 0, 1463827927039889790, Guage, h2b("00 00 00 00 00 43 cf 41")}
var testDate = time.Date(2013, time.March, 14, 21, 19, 53, 804828672, time.UTC)

func TestValueTime(t *testing.T) {
	result := testValue.Time()
	if !result.Equal(testDate) {
		t.Errorf("expected %v, got %v", testDate, result)
	}
}

func TestValueSeconds(t *testing.T) {
	result := testValue.TimeUnix()
	expected := testDate.Unix()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestValueNanoSeconds(t *testing.T) {
	result := testValue.TimeUnixNano()
	expected := testDate.UnixNano()
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
