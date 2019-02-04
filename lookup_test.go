package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupIP(t *testing.T) {
	prefixes := []Prefix{
		{
			IPPrefix: "192.168.0.0/24", Region: "eu-west-1", Service: "EC2",
		},
		{
			IPPrefix: "192.168.1.0/24", Region: "eu-west-2", Service: "S3",
		},
		{
			IPPrefix: "192.168.2.0/24", Region: "eu-west-3", Service: "CLOUDFRONT",
		},
	}
	IPv6Prefixes := []IPv6Prefix{
		{
			IPv6Prefix: "2001:800::/21", Region: "us-east-2", Service: "GLOBALACCELERATOR",
		},
		{
			IPv6Prefix: "2600:1fa0:4000::/40", Region: "us-east-1", Service: "CODEBUILD",
		},
	}
	doc := IPRangeDoc{
		Prefixes:     prefixes,
		IPv6Prefixes: IPv6Prefixes,
	}

	lpi := lookupPrefixInput{
		doc: doc,
		ip:  "192.168.0.100",
	}
	lpo, err := lookupPrefix(lpi)
	assert.NoError(t, err)
	assert.Equal(t, lpo.doc.Prefixes[0].IPPrefix, "192.168.0.0/24")
	assert.Equal(t, lpo.doc.Prefixes[0].Region, "eu-west-1")
	assert.Equal(t, lpo.doc.Prefixes[0].Service, "EC2")

	l6pi := lookupPrefixInput{
		doc: doc,
		ip:  "2600:1fa0:4060:c11:34da:e8e9::",
	}
	lpo, err = lookupPrefix(l6pi)
	assert.NoError(t, err)
	assert.Equal(t, lpo.doc.IPv6Prefixes[0].IPv6Prefix, "2600:1fa0:4000::/40")
	assert.Equal(t, lpo.doc.IPv6Prefixes[0].Region, "us-east-1")
	assert.Equal(t, lpo.doc.IPv6Prefixes[0].Service, "CODEBUILD")
}
