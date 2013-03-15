package gocollectd

import (
	"testing"
	"fmt"
	"strings"
	"strconv"
	"reflect"
)

func hexdumpToBytes(s ...string) ([]byte, error) {
	joined := strings.Join(s, " ")
	a := strings.Split(joined, " ")
	b := make([]byte, len(a))
	for i, n := range(a) {
		byt, err := strconv.ParseUint(n, 16, 8)
		if err != nil {
			return []byte{}, err
		}
		b[i] = byte(byt)
	}
	return b, nil
}

func h2b(s ...string) []byte {
	b, err := hexdumpToBytes(s...)
	if err != nil {
	   panic(fmt.Sprintf("failed to convert %v to bytes", s))
	}
	return b
}

func TestParse5(t *testing.T) {
	// unfortunately a real hexdump seems the best way to test this
	b := h2b(
		"00 00 00 0f 6c 61 70 74 6f 70 2e 6c 61 6e 00", // hostname: "laptop.lan"
		"00 08 00 0c 14 50 8f be 73 82 51 7e",          // time, hi res
		"00 09 00 0c 00 00 00 02 80 00 00 00",          // interval
		"00 02 00 0b 6d 65 6d 6f 72 79 00",             // plugin: memory
		"00 05 00 0a 77 69 72 65 64 00",                // type instance: wired
		"00 06 00 0f 00 01 01 00 00 00 00 00 43 cf 41", // value
		"00 08 00 0c 14 50 8f be 73 82 94 9a",          // time, hi res
		"00 02 00 0e 69 6e 74 65 72 66 61 63 65 00",    // plugin: interface
		"00 03 00 08 6c 6f 30 00",                      // instance: lo0
		"00 04 00 0e 69 66 5f 6f 63 74 65 74 73 00",    // type: if_octets
		"00 05 00 05 00",                               // type instance: nil
		"00 06 00 18 00 02 02 02 00 00 00 00 00 88 07 8b 00 00 00 00 00 88 07 8c", // 2 more values, note: the second one was manipulated to check order
		"00 08 00 0c 14 50 8f be 73 84 40 6c",          // a new time
		"00 04 00 0f 69 66 5f 70 61 63 6b 65 74 73 00", // plugin: ifpackets
		"00 06 00 18 00 02 02 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00", // 2 more values
	)
	expected := []Value{
		{"laptop.lan", "memory", "", "", "wired", 0, 1463827927039889790, 1, h2b("00 00 00 00 00 43 cf 41")},
		{"laptop.lan", "interface", "lo0", "if_octets", "", 0, 1463827927039906970, 2, h2b("00 00 00 00 00 88 07 8b")},
		{"laptop.lan", "interface", "lo0", "if_octets", "", 1, 1463827927039906970, 2, h2b("00 00 00 00 00 88 07 8c")},
		{"laptop.lan", "interface", "lo0", "if_packets", "", 0, 1463827927040016492, 2, h2b("00 00 00 00 00 00 00 00")},
		{"laptop.lan", "interface", "lo0", "if_packets", "", 1, 1463827927040016492, 2, h2b("00 00 00 00 00 00 00 00")},
	}
	result := Parse(b)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v\n", expected, result)
	}
}

func TestParse4(t *testing.T) {
	// worse than a hexdump: this is an edited hexdump based on the spec
	// as I don't have a copy of collectd4 running right now
	b := h2b(
		"00 00 00 0f 6c 61 70 74 6f 70 2e 6c 61 6e 00", // hostname: "laptop.lan"
		"00 01 00 0c 00 00 00 00 51 42 3E F9",          // time, low res
		"00 07 00 0c 00 00 00 00 00 00 00 0A",          // interval
		"00 02 00 0b 6d 65 6d 6f 72 79 00",             // plugin: memory
		"00 05 00 0a 77 69 72 65 64 00",                // type instance: wired
		"00 06 00 0f 00 01 01 00 00 00 00 00 43 cf 41", // value
		"00 01 00 0c 00 00 00 00 51 42 3E FA",          // time, low res
		"00 02 00 0e 69 6e 74 65 72 66 61 63 65 00",    // plugin: interface
		"00 03 00 08 6c 6f 30 00",                      // instance: lo0
		"00 04 00 0e 69 66 5f 6f 63 74 65 74 73 00",    // type: if_octets
		"00 05 00 05 00",                               // type instance: nil
		"00 06 00 18 00 02 02 02 00 00 00 00 00 88 07 8b 00 00 00 00 00 88 07 8c", // 2 more values, note: the second one was manipulated to check order
		"00 04 00 0f 69 66 5f 70 61 63 6b 65 74 73 00", // plugin: ifpackets
		"00 06 00 18 00 02 02 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00", // 2 more values
	)

	expected := []Value{
		{"laptop.lan", "memory", "", "", "wired", 0, 1463827926175711232, 1, h2b("00 00 00 00 00 43 cf 41")},
		{"laptop.lan", "interface", "lo0", "if_octets", "", 0, 1463827927249453056, 2, h2b("00 00 00 00 00 88 07 8b")},
		{"laptop.lan", "interface", "lo0", "if_octets", "", 1, 1463827927249453056, 2, h2b("00 00 00 00 00 88 07 8c")},
		{"laptop.lan", "interface", "lo0", "if_packets", "", 0, 1463827927249453056, 2, h2b("00 00 00 00 00 00 00 00")},
		{"laptop.lan", "interface", "lo0", "if_packets", "", 1, 1463827927249453056, 2, h2b("00 00 00 00 00 00 00 00")},
	}
	result := Parse(b)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected\n%v\ngot\n%v\n", expected, result)
	}
}


func TestParseEmpty(t *testing.T) {
	result := Parse([]byte{})
	if !reflect.DeepEqual(result, []Value{}) {
		t.Errorf("expected [] got %v", result)
	}
}