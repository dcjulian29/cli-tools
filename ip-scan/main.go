package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
)

type EndpointInfo struct {
	ip       net.IP
	fqdn     string
	ping     time.Duration
	ttl      int
	attempts int
}

func main() {
	ips, err := expand_cidr(os.Args[1])
	if err != nil {
		fmt.Println("Error scanning IPs:", err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan EndpointInfo)
	sortedResults := make([]EndpointInfo, 0)

	for _, ip := range ips {
		wg.Add(1)

		go func(ip string) {
			defer wg.Done()

			attempts := 1
			pingTime := time.Duration(-1)
			ttl := 0

			for i := 0; i < 3; i++ {
				pingTime, ttl = pingIP(ip)

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
					ttl:      ttl,
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
	table.SetHeader([]string{"IP Address", "FQDN", "Ping Time", "TTL"})
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
			fmt.Sprintf("%d", info.ttl),
		})
	}

	fmt.Println("")

	table.Render()

	fmt.Printf("\n %d host(s) responded\n", len(sortedResults))
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

func pingIP(ip string) (time.Duration, int) {
	addr := net.IPAddr{IP: net.ParseIP(ip)}
	icmpPacket := []byte{8, 0, 247, 255, 0, 0, 0, 0}
	conn, err := net.DialIP("ip4:icmp", nil, &addr)
	if err != nil {
		return time.Duration(-1), 0
	}

	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return time.Duration(-1), 0
	}

	start := time.Now()

	_, err = conn.Write(icmpPacket)
	if err != nil {
		return time.Duration(-1), 0
	}

	reply := make([]byte, 1500)

	_, err = conn.Read(reply)
	if err != nil {
		return time.Duration(-1), 0
	}

	return time.Since(start), int(reply[8])
}
