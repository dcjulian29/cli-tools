package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]

	dst, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("Ping"),
		},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	reply := make([]byte, 1500)
	start := time.Now()

	if _, err = conn.WriteTo(msgBytes, dst); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		if err = conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			log.Fatal(err)
		}

		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			log.Fatal(err)
		}

		duration := time.Since(start)

		if msg, err = icmp.ParseMessage(protocolICMP, reply[:n]); err != nil {
			log.Fatal(err)
		}

		switch msg.Type {
		case ipv4.ICMPTypeEchoReply:
			echoReply, ok := msg.Body.(*icmp.Echo)
			if !ok {
				log.Fatal("invalid ICMP Echo Reply message")
				return
			}

			if peer.String() == host && echoReply.ID == os.Getpid()&0xffff && echoReply.Seq == 1 {
				fmt.Printf("reply from %s: time=%v\n", dst.String(), duration)
				return
			}

		default:
			fmt.Printf("unexpected ICMP message type: %v\n", msg.Type)
		}
	}
}
