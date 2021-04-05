package main

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type outputInput struct {
	Doc           IPRangeDoc
	fourOnly      bool
	sixOnly       bool
	encoding      string
	textSeparator string
	fields        string
}

type outItems struct {
	SyncToken    string       `json:"syncToken" yaml:"syncToken"`
	CreateDate   string       `json:"createDate" yaml:"createDate"`
	Prefixes     []Prefix     `json:"prefixes" yaml:"prefixes"`
	IPv6Prefixes []IPv6Prefix `json:"ipv6_prefixes" yaml:"ipv6_prefixes"`
}

func output(input outputInput) string {
	var strBuilder strings.Builder
	var oi outItems
	oi.SyncToken = input.Doc.SyncToken
	oi.CreateDate = input.Doc.CreateDate
	oi.Prefixes = input.Doc.Prefixes
	oi.IPv6Prefixes = input.Doc.IPv6Prefixes
	switch input.encoding {
	case strText:
		fields := strings.Split(input.fields, ",")
		if !input.sixOnly {
			for _, p := range input.Doc.Prefixes {
				if stringInSlice("service", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(fmt.Sprintf("%s%s", p.Service, input.textSeparator))
				}
				if stringInSlice("region", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(fmt.Sprintf("%s%s", p.Region, input.textSeparator))
				}
				if stringInSlice("cidr", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(p.IPPrefix)
				}
				strBuilder.WriteString("\n")
			}
		}
		if !input.fourOnly {
			for _, p := range input.Doc.IPv6Prefixes {
				if stringInSlice("service", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(fmt.Sprintf("%s%s", p.Service, input.textSeparator))
				}
				if stringInSlice("region", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(fmt.Sprintf("%s%s", p.Region, input.textSeparator))
				}
				if stringInSlice("cidr", fields) || stringInSlice("all", fields) {
					strBuilder.WriteString(p.IPv6Prefix)
				}
				strBuilder.WriteString("\n")
			}
		}
	case "json":
		o, _ := json.MarshalIndent(oi, "", "  ")
		strBuilder.WriteString(string(o))
	case "yaml":
		o, _ := yaml.Marshal(oi)
		strBuilder.WriteString(string(o))
	}
	return strBuilder.String()
}
