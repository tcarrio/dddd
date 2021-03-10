[![Go Report Card](https://goreportcard.com/badge/github.com/tcarrio/dddd)](https://goreportcard.com/report/github.com/tcarrio/dddd) [![builds.sr.ht status](https://builds.sr.ht/~tcarrio/dddd/commits/.build.yml.svg)](https://builds.sr.ht/~tcarrio/dddd/commits/.build.yml?)

# dddd

Self-hostable DynDNS Alternative

## Overview

Having looked at other dynamic DNS solutions and realizing how simple the idea was, I decided to throw together a CLI that would update a DNS entry from a running server. The initial motivation came to me while setting up my at home VPN and looking up my IP address on DuckDuckGo. It returned an IP address via the Answers API, so my solution unironically uses the same. **dddd** hits the DuckDuckGo Answers API, parses the returned IP address, and will create or update a DNS entry in Cloudflare. I believe this solution would be easy to throw in a cronjob for scheduled updates, as well as manual invocation. Currently it is heavily driven by arguments or environment variables.

## Usage

```
NAME:
   dddd - perform dynamic DNS updates in Cloudflare using DuckDuckGo Answers

USAGE:
   dddd [global options] command [command options] [arguments...]

VERSION:
   0.2.4

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-key key, -k key        Cloudflare API key [$CF_API_KEY]
   --email email, -e email      Cloudflare email [$CF_API_EMAIL]
   --ip IP, -i IP               IP address for the DNS entry [$DNS_IP]
   --name name, -n name         The name for the DNS entry (default: "@") [$DNS_NAME]
   --domain domain, -d domain   The domain to modify in Cloudflare [$DNS_DOMAIN]
   --type type, -t type         The type of the record [$DNS_TYPE]
   --log-level value, -l value  Set the log level of logger [$LOG_LEVEL]
   --help, -h                   show help
   --version, -v                print the version
```
