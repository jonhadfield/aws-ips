package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	ipURL     = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	strCIDR   = "cidr"
	strPrefix = "prefix"
	strJSON   = "json"
	strYAML   = "yaml"
	strText   = "text"
)

var (
	validFields    = []string{"cidr", "service", "region", "all"}
	validEncodings = []string{strJSON, strYAML, strText}
)

type Prefix struct {
	IPPrefix string `json:"ip_prefix" yaml:"ip_prefix"`
	Region   string `json:"region" yaml:"region"`
	Service  string `json:"service" yaml:"service"`
}

type IPv6Prefix struct {
	IPv6Prefix string `json:"ipv6_prefix" yaml:"ipv6_prefix"`
	Region     string `json:"region" yaml:"region"`
	Service    string `json:"service" yaml:"service"`
}

type IPRangeDoc struct {
	SyncToken    string       `json:"syncToken"`
	CreateDate   string       `json:"createDate"`
	Prefixes     []Prefix     `json:"prefixes" yaml:"prefixes"`
	IPv6Prefixes []IPv6Prefix `json:"ipv6_prefixes" yaml:"ipv6_prefixes"`
}

// overwritten at build time
var version, versionOutput, tag, sha, buildDate string

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	msg, display, err := startCLI(os.Args)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	if display && msg != "" {
		fmt.Print(msg)
	}
	os.Exit(0)
}

func startCLI(args []string) (msg string, display bool, err error) {
	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Name = "aws-ips"
	app.Version = versionOutput
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Jon Hadfield",
			Email: "jon@lessknown.co.uk",
		},
	}
	app.HelpName = "-"
	app.Usage = "aws-ips"
	app.Description = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region",
			Usage: "AWS region",
		},
		cli.StringFlag{
			Name:  "service",
			Usage: "AWS service",
		},
		cli.StringFlag{
			Name:  "ip",
			Usage: "find prefix by ip",
		},
		cli.BoolFlag{
			Name:  "4only",
			Usage: "only output IPv4",
		},
		cli.BoolFlag{
			Name:  "6only",
			Usage: "only output IPv6",
		},
		cli.StringFlag{
			Name:  "fields",
			Value: "all",
			Usage: "cidr, region, service, all",
		},
		cli.StringFlag{
			Name:  "encoding",
			Value: "text",
			Usage: "output encoding (json, yaml, text)",
		},
		cli.StringFlag{
			Name:  "separator",
			Value: " | ",
			Usage: "field separator for text encoding",
		},
		cli.StringFlag{
			Name:  "from-file",
			Usage: "file downloaded from https://ip-ranges.amazonaws.com/ip-ranges.json",
		},
		cli.BoolFlag{
			Name:  "silent",
			Usage: "no stdout",
		},
	}

	app.Action = func(c *cli.Context) error {
		region := c.String("region")
		service := c.String("service")
		ip := c.String("ip")
		encoding := c.String("encoding")
		fields := c.String("fields")
		textSeparator := c.String("separator")
		if c.Bool("silent") {
			display = false
		} else {
			display = true
		}
		if ip != "" {
			if region != "" || service != "" {
				_, _ = fmt.Fprintf(c.App.Writer, "error: ip cannot be used with region or service.\n")
				os.Exit(1)
			}
			// check ip valid
			if strings.HasSuffix(ip, "/32") {
				ip = ip[:len(ip)-3]
			}
			if net.ParseIP(ip) == nil {
				_, _ = fmt.Fprintf(c.App.Writer, "error: ip is invalid.\n")
				os.Exit(1)
			}
		}

		for _, f := range strings.Split(fields, ",") {
			if !stringInSlice(f, validFields) {
				_, _ = fmt.Fprintf(c.App.Writer, "error: fields must be one of: %s.\n",
					strings.Join(validFields, ", "))
				os.Exit(1)
			}
		}

		if encoding != "" && !stringInSlice(encoding, validEncodings) {
			_, _ = fmt.Fprintf(c.App.Writer, "error: encoding must be one of: %s.\n",
				strings.Join(validEncodings, ", "))
			os.Exit(1)
		}

		loaded, err := loadRanges(c.String("from-file"))
		if err != nil {
			return err
		}
		outputDoc := loaded
		switch {
		case ip != "":
			var lookupPrefixInput lookupPrefixInput
			lookupPrefixInput.doc = loaded
			lookupPrefixInput.ip = ip
			var lookupOutput lookupPrefixOutput
			lookupOutput, err = lookupPrefix(lookupPrefixInput)
			if err != nil {
				return err
			}
			outputDoc = lookupOutput.doc
			// default to outputting all fields
			if fields == "" {
				fields = "all"
			}
		case service != "" || region != "":
			var frInput filterRangesInput
			frInput.Doc = loaded
			frInput.region = region
			frInput.service = service
			var filteredOutput filterRangesOutput
			filteredOutput = filterRanges(frInput)
			outputDoc = filteredOutput.Doc
			if err != nil {
				return err
			}
			// default to outputting CIDR
			if fields == "" {
				fields = strCIDR
			}
		default:
			// default to returning all, unfiltered
			if fields == "" {
				fields = strPrefix
			}
		}

		// only output text format if there are entries to output
		if encoding == strText && (len(outputDoc.Prefixes) == 0 && len(outputDoc.IPv6Prefixes) == 0) {
			return nil
		}
		msg = output(outputInput{
			fourOnly:      c.Bool("4only"),
			sixOnly:       c.Bool("6only"),
			Doc:           outputDoc,
			encoding:      encoding,
			textSeparator: textSeparator,
			fields:        fields,
		})

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	return msg, display, app.Run(args)
}
