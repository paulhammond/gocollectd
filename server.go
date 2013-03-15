package gocollectd

import (
	"log"
	"net"
)

func Listen(addr string, c chan Packet) {
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
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println("error: Failed to recieve packet", err)
		} else {
			packets, err := Parse(buf[0:n])
			if err != nil {
				log.Println("error: Failed to recieve packet", err)
			}
			for _, p := range *packets {
				c <- p
			}
		}
	}
}
