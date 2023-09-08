package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/theckman/yacspin"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type EndpointInfo struct {
	ip       net.IP
	fqdn     string
	ping     time.Duration
	attempts int
}

func main() {
	start := time.Now()
	ips, err := expand_cidr(os.Args[1])
	if err != nil {
		fmt.Println("Error scanning IPs:", err)
		return
	}

	spinner, _ := yacspin.New(yacspin.Config{
		Frequency: 100 * time.Millisecond,
		Colors:    []string{"fgYellow"},
		CharSet:   yacspin.CharSets[69],
	})

	spinner.Start()

	var wg sync.WaitGroup
	results := make(chan EndpointInfo)
	sortedResults := make([]EndpointInfo, 0)

	for _, ip := range ips {
		wg.Add(1)

		go func(ip string) {
			defer wg.Done()

			attempts := 1
			pingTime := time.Duration(-1)

			for i := 0; i < 3; i++ {
				pingTime = pingIP(ip)

				if pingTime > 0 {
					break
				}

				attempts++
			}

			if pingTime > 0 {
				info := EndpointInfo{
					ip:       net.ParseIP(ip),
					fqdn:     getFQDN(ip),
					ping:     pingTime,
					attempts: attempts,
				}

				results <- info
			}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for item := range results {
		sortedResults = append(sortedResults, item)
	}

	sort.Slice(sortedResults, func(i, j int) bool {
		return bytes.Compare(sortedResults[i].ip, sortedResults[j].ip) < 0
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP Address", "FQDN", "Ping Time"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t\t")
	table.SetNoWhiteSpace(true)

	for _, info := range sortedResults {
		table.Append([]string{
			info.ip.String(),
			info.fqdn,
			info.ping.String(),
		})
	}

	spinner.Stop()

	fmt.Println("")

	table.Render()

	duration := time.Since(start).Truncate(time.Millisecond)
	fmt.Printf("\n  %d host(s) responded in %v\n", len(sortedResults), duration)
}

func expand_cidr(ipRange string) ([]string, error) {
	ips := []string{}

	ip, ipNet, err := net.ParseCIDR(ipRange)
	if err != nil {
		return ips, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func getFQDN(ip string) string {
	fqdn, err := net.LookupAddr(ip)
	if err != nil {
		return ""
	}

	return fqdn[0]
}

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
