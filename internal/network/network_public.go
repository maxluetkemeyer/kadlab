package network

import (
	"context"
	"fmt"
	"net"
)

type PublicNetwork struct{}

func (network *PublicNetwork) ResolveDNS(ctx context.Context, domain string) ([]string, error) {
	// Perform a DNS lookup for the given domain
	ipAddresses, err := net.LookupIP(domain)

	if err != nil {
		err = fmt.Errorf("failed to resolve domain '%s': %s", domain, err)
		return nil, err
	}

	// Convert net.IP addresses to strings and return them
	var ips []string

	for _, ipAddress := range ipAddresses {
		ips = append(ips, ipAddress.String())
	}

	return ips, nil
}
