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
	"github.com/theckman/yacspin"
)

type EndpointInfo struct {
	ip   net.IP
	fqdn string
}

func main() {
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

			fqdn := getFQDN(ip)

			if len(fqdn) > 0 {
				info := EndpointInfo{
					ip:   net.ParseIP(ip),
					fqdn: fqdn,
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
	table.SetHeader([]string{"IP Address", "FQDN"})
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
		})
	}

	spinner.Stop()

	fmt.Println("")

	table.Render()
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
