package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSONFile(t *testing.T) {
	ranges, err := loadFromFile("testdata/ip-ranges.json")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(ranges.Prefixes))
	assert.NotEqual(t, 0, len(ranges.IPv6Prefixes))
	assert.NotEmpty(t, ranges.SyncToken)
	assert.NotEmpty(t, ranges.CreateDate)
}

func TestLoadURL(t *testing.T) {
	ranges, err := loadFromURL()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(ranges.Prefixes))
	assert.NotEqual(t, 0, len(ranges.IPv6Prefixes))
	assert.NotEmpty(t, ranges.SyncToken)
	assert.NotEmpty(t, ranges.CreateDate)
}
