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
		fmt.Println("-")
		for i, name := range packet.ValueNames() {
			values, _ := packet.ValueNumbers()
			fmt.Printf("%s %35s %v\n", packet.Time().Format(time.RFC3339), name, values[i])
		}
	}
}
