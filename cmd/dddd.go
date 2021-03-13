package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/tcarrio/dddd/pkg/dddd"
)

const (
	// Version is the version of dddd tool
	Version = "0.2.5"
)

func main() {
	app := cli.NewApp()
	app.Name = "dddd"
	app.Usage = "perform dynamic DNS updates in Cloudflare using DuckDuckGo Answers"
	app.Version = Version
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key, k",
			EnvVar: "CF_API_KEY",
			Usage:  "Cloudflare API `key`",
		},
		cli.StringFlag{
			Name:   "email, e",
			EnvVar: "CF_API_EMAIL",
			Usage:  "Cloudflare `email`",
		},
		cli.StringFlag{
			Name:   "ip, i",
			EnvVar: "DNS_IP",
			Usage:  "`IP` address for the DNS entry",
		},
		cli.StringFlag{
			Name:   "name, n",
			EnvVar: "DNS_NAME",
			Usage:  "The `name` for the DNS entry",
			Value:  "@",
		},
		cli.StringFlag{
			Name:   "domain, d",
			EnvVar: "DNS_DOMAIN",
			Usage:  "The `domain` to modify in Cloudflare",
		},
		cli.StringFlag{
			Name:   "type, t",
			EnvVar: "DNS_TYPE",
			Usage:  "The `type` of the record",
		},
		cli.StringFlag{
			Name:   "log-level, l",
			EnvVar: "LOG_LEVEL",
			Usage:  "Set the log level of logger",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Use the specified config file",
			EnvVar: "DDDD_CONFIG",
		},
		cli.BoolFlag{
			Name:   "make, m",
			Usage:  "Shows the version for make",
			Hidden: true,
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	if c.Bool("make") {
		fmt.Println(Version)
		return
	}

	args := dddd.Config{}

	configLocation := c.String("config")
	if len(configLocation) > 0 {
		configArgs, err := dddd.ParseConfig(configLocation)
		if err != nil {
			fmt.Printf("Skipping configuration file at %s\n", configLocation)
		} else {
			dddd.MergeConfig(&args, configArgs)
		}

	}

	dddd.MergeConfig(&args, dddd.Config{
		Key:        c.String("api-key"),
		Email:      c.String("email"),
		Name:       c.String("name"),
		Domain:     c.String("domain"),
		IP:         c.String("ip"),
		RecordType: c.String("type"),
		LogLevel:   c.String("log-level"),
	})

	dddd.Run(args)
}
