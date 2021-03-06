package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
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
var (
	version, versionOutput, tag, sha, buildDate string
	msg                                         string
	rCode                                       int
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	display, err := startCLI(os.Args)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
	if display && msg != "" {
		fmt.Print(msg)
	}
	os.Exit(rCode)
}

func startCLI(args []string) (display bool, err error) {
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
			Name:  "region, r",
			Usage: "AWS region",
		},
		cli.StringFlag{
			Name:  "service, s",
			Usage: "AWS service",
		},
		cli.StringFlag{
			Name:  "ip",
			Usage: "find prefix by ip",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "find prefix by name",
		},
		cli.BoolFlag{
			Name:  "4only, 4",
			Usage: "only output IPv4",
		},
		cli.BoolFlag{
			Name:  "6only, 6",
			Usage: "only output IPv6",
		},
		cli.BoolFlag{
			Name:  "no-amazon, na",
			Usage: "exclude matches with service AMAZON",
		},
		cli.StringFlag{
			Name:  "fields, f",
			Value: "all",
			Usage: "cidr, region, service, all",
		},
		cli.StringFlag{
			Name:  "encoding, e",
			Value: "text",
			Usage: "output encoding (json, yaml, text)",
		},
		cli.StringFlag{
			Name:  "separator, sep",
			Value: " | ",
			Usage: "field separator for text encoding",
		},
		cli.StringFlag{
			Name:  "from-file, file",
			Usage: "file downloaded from https://ip-ranges.amazonaws.com/ip-ranges.json",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "no stdout",
		},
	}

	app.Action = func(c *cli.Context) error {
		region := c.String("region")
		service := c.String("service")
		ip := c.String("ip")
		name := c.String("name")
		noAmazon := c.Bool("no-amazon")
		encoding := c.String("encoding")
		fields := c.String("fields")
		textSeparator := c.String("separator")
		if c.Bool("quiet") {
			display = false
		} else {
			display = true
		}
		if ip != "" && name != "" {
			_, _ = fmt.Fprintf(c.App.Writer, "error: ip cannot be used in combination with name.\n")
			os.Exit(1)
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
		if name != "" {
			if region != "" || service != "" {
				_, _ = fmt.Fprintf(c.App.Writer, "error: name cannot be used with region or service.\n")
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
		case name != "":
			// Get IP for name
			ips, err := net.LookupHost(name)
			if err != nil {
				_, _ = fmt.Fprintf(c.App.Writer, "error: failed to resolve: %s\n", name)
				os.Exit(1)
			}
			ip = ips[0]
			var lookupPrefixInput lookupPrefixInput
			lookupPrefixInput.doc = loaded
			lookupPrefixInput.ip = ip
			lookupPrefixInput.noAmazon = noAmazon
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
			frInput.noAmazon = noAmazon
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

		if len(outputDoc.Prefixes) == 0 && len(outputDoc.IPv6Prefixes) == 0 {
			rCode = 2
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
	return display, app.Run(args)
}
