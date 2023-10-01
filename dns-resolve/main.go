package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]

	dst, err := net.ResolveIPAddr("ip", host)

	if err != nil {
		fmt.Printf("%v", err)
	} else {
		fmt.Printf("lookup %s: %v", host, dst)
	}
}
