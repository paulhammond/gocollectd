// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package main

import (
	"fmt"
	collectd "github.com/paulhammond/gocollectd"
	"time"
)

func main() {
	c := make(chan collectd.Packet)
	go collectd.Listen("127.0.0.1:25827", c)

	for {
		packet := <-c
		for i, name := range packet.ValueNames() {
			values, _ := packet.ValueNumbers()
			newpacket := ""
			if (i == 0) {
				newpacket = "."
			}
			fmt.Printf("%2s %s %35s %v\n", newpacket, packet.Time().Format(time.RFC3339), name, values[i])
		}
	}
}
