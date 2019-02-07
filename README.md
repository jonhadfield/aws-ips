# aws-ips
check network addresses for AWS Services

[![Build Status](https://www.travis-ci.org/jonhadfield/aws-ips.svg?branch=master)](https://www.travis-ci.org/jonhadfield/aws-ips) [![Go Report Card](https://goreportcard.com/badge/github.com/jonhadfield/aws-ips)](https://goreportcard.com/report/github.com/jonhadfield/aws-ips)

## about

AWS [publish a list](https://ip-ranges.amazonaws.com/ip-ranges.json) of IP ranges that correspond to their services and regions. Linked to from [here](https://docs.aws.amazon.com/general/latest/gr/aws-ip-ranges.html).  
aws-ips is a utility that downloads the latest list (or use a local copy of the list) and enables you to:

- list and filter prefixes (by region, service, cidr)
- reverse lookup of IP address or FQDN to AWS range
- output results as JSON, YAML, or Text

## changelog
0.0.1 - initial

## installation
Download the latest release here: https://github.com/jonhadfield/aws-ips/releases

#### macOS and Linux
  
Install:  
``
$ install <aws-ips binary> /usr/local/bin/aws-ips
``  
#### Windows
  
An installer is planned, but for now...  
Download the binary 'aws-ips_windows_amd64.exe' and rename to aws-ips.exe

### example usage  

#### help
``
aws-ips --help
``  

#### filtering by service and region

``
aws-ips --service s3 --region eu-west-1
``  

#### reverse lookup by IP

``
aws-ips --ip 54.208.0.0
``

#### AWS service relating to a domain name

``
aws-ips --name www.netflix.com
``
