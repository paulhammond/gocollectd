// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package gocollectd

import (
	"log"
	"net"
)

// Listen creates a UDP server that parses collectd data into packets and
// sends them over a channel.
// The caller is responsible for ensuring that the channel is ready to receive
// data, otherwise packets will be dropped
func Listen(addr string, c chan<- Packet) {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("fatal: failed to resolve address", err)
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatalln("fatal: failed to listen", err)
	}
	for {
		buf := make([]byte, 1452)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("error: Failed to receive packet", err)
		} else {
			packets, err := Parse(buf[0:n])
			if err != nil {
				log.Println("error: Failed to receive packet", err)
			}
			for _, p := range *packets {
				select {
				case c <- p:
					// packet sent to channel
				default:
					// don't block if channel isn't ready
					log.Println("error: Channel not ready for write. Packet dropped.")
				}
			}
		}
	}
}
