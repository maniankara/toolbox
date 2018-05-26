package main

import (
	"fmt"
	"log"
	"net"
)

// Validates ipv4
// Ref: https://gist.github.com/maniankara/f321a15a9bb4c9e4e92b2829d7d2f169
func isValidIpv4(given net.IP) bool {
	var v = true
	if given.To4() == nil {
		v = false
	}

	// Skip loop back
	if given.IsLoopback() {
		v = false
	}

	return v
}

// Returns a list of ip(s)
// Ref: https://gist.github.com/maniankara/f321a15a9bb4c9e4e92b2829d7d2f169
func getMyIPs() []net.IP {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("Error obtaining address from interface: %s continuing...", i)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// Skip the loop back and other garbage
			if !isValidIpv4(ip) {
				continue
			}
			ips = append(ips, ip)
		}
	}
	return ips
}

// increments ip and sees it does not resolve in the network
// until the given count is reached
func findFreeIPs(ip net.IP, count int) []net.IP {
	var ips []net.IP
	ip4 := ip.To4()
	counter := 0
	for {
		ip4[3]++
		dest, err := net.LookupAddr(ip4.String())
		if len(dest) == 0 && err != nil {
			ips = append(ips, ip4.To16())
		}
		counter++
		if counter == count {
			break
		}
	}
	return ips
}

func main() {
	// get ips from your local interface
	ips := getMyIPs()
	// Start guessing
	for _, ip := range ips {
		fmt.Println(findFreeIPs(ip, 2))
	}

}
