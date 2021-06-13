package utils

import (
	"net"
)

// Based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
func GetBroadcastAddress(ipnet *net.IPNet) (string, error) {
	ips := []string{}
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// The last address is the broadcast address
	return ips[len(ips)-1], nil
}

// See http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
