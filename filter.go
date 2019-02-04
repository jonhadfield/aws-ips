package main

import "strings"

func filterRanges(input filterRangesInput) filterRangesOutput {
	var output filterRangesOutput
	output.Doc.SyncToken = input.Doc.SyncToken
	output.Doc.CreateDate = input.Doc.CreateDate
	for _, r := range input.Doc.Prefixes {
		var match bool
		if input.region != "" {
			if strings.EqualFold(r.Region, input.region) {
				match = true
			} else {
				continue
			}
		}
		if input.service != "" {
			if strings.EqualFold(r.Service, input.service) {
				match = true
			} else {
				continue
			}
		}
		if match {
			output.Doc.Prefixes = append(output.Doc.Prefixes, r)
		}
	}
	for _, ip6r := range input.Doc.IPv6Prefixes {
		var match bool
		if input.region != "" {
			if strings.EqualFold(ip6r.Region, input.region) {
				match = true
			} else {
				continue
			}
		}

		if input.service != "" {
			if strings.EqualFold(ip6r.Service, input.service) {
				match = true
			} else {
				continue
			}
		}
		if match {
			output.Doc.IPv6Prefixes = append(output.Doc.IPv6Prefixes, ip6r)
		}
	}
	return output
}

type filterRangesInput struct {
	Doc     IPRangeDoc
	region  string
	service string
}

type filterRangesOutput struct {
	Doc IPRangeDoc
}
