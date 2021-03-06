package main

import (
	"net"
	"strings"
)

func lookupPrefix(input lookupPrefixInput) (output lookupPrefixOutput, err error) {
	output.doc.CreateDate = input.doc.CreateDate
	output.doc.SyncToken = input.doc.SyncToken
	netIP := net.ParseIP(input.ip)
	if netIP != nil {
		if strings.Contains(netIP.String(), ":") {
			for _, r := range input.doc.IPv6Prefixes {
				if input.noAmazon && strings.ToLower(r.Service) == "amazon" {
					continue
				}
				var netPrefix *net.IPNet
				_, netPrefix, err = net.ParseCIDR(r.IPv6Prefix)
				if netPrefix.Contains(netIP) {
					output.doc.IPv6Prefixes = append(output.doc.IPv6Prefixes, r)
				}
			}
		} else {
			for _, r := range input.doc.Prefixes {
				if input.noAmazon && strings.ToLower(r.Service) == "amazon" {
					continue
				}
				var netPrefix *net.IPNet
				_, netPrefix, err = net.ParseCIDR(r.IPPrefix)
				if netPrefix.Contains(netIP) {
					output.doc.Prefixes = append(output.doc.Prefixes, r)
				}
			}
		}
	}

	return output, err
}

type lookupPrefixInput struct {
	doc      IPRangeDoc
	ip       string
	noAmazon bool
}

type lookupPrefixOutput struct {
	doc IPRangeDoc
}
