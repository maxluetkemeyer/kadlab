package network

import (
	"context"
	"log"
	"net"
)

type PublicNetwork struct{}

func (network *PublicNetwork) ResolveDNS(ctx context.Context, domain string) []string {
	// Perform a DNS lookup for the given domain
	ipAddresses, err := net.LookupIP(domain)

	if err != nil {
		log.Fatalf("Failed to resolve domain '%s': %s", domain, err)
		return []string{}
	}

	// Convert net.IP addresses to strings and return them
	var ips []string

	for _, ipAddress := range ipAddresses {
		ips = append(ips, ipAddress.String())
	}

	return ips
}
