//go:build linux || (darwin && !cgo)

/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func pingIP(ip string) time.Duration {
	addr := net.UDPAddr{IP: net.ParseIP(ip)}
	conn, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  0,
			Data: []byte("ping"),
		},
	}

	wb, _ := msg.Marshal(nil)
	start := time.Now()
	conn.WriteTo(wb, &addr)
	rb := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, peer, err := conn.ReadFrom(rb)

	if err == nil {
		duration := time.Since(start)
		rm, err := icmp.ParseMessage(1, rb[:n])
		if err == nil {
			if rm.Type == ipv4.ICMPTypeEchoReply {
				echoReply, ok := msg.Body.(*icmp.Echo)
				if ok {
					if peer.(*net.UDPAddr).IP.String() == ip && echoReply.Seq == 0 {
						return duration
					}
				}
			}
		}
	}

	return 0
}
