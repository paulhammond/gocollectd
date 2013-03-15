package main

import(
	collectd "github.com/paulhammond/gocollectd"
	"fmt"
)

func main() {
	c := make(chan collectd.Packet)
	go collectd.Listen("127.0.0.1:25827", c)

	for {
		packet := <- c
		fmt.Println(packet)
	}

}

