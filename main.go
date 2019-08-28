package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/cloudflare/cloudflare-go"
	"github.com/urfave/cli"
)

const (
	exampleAnswer = "Your IP address is 68.32.2.167 in <a href=\"https://duckduckgo.com/?q=Royal%20Oak%2C%20Michigan%2C%20United%20States%20(48073)&iar=maps_maps\">Royal Oak, Michigan, United States (48073)</a>"
	answerAPI     = "https://api.duckduckgo.com/?q=what%27s+my+ip&format=json"
)

// Answer is the JSON response for the DuckDuckGo Answer API
type Answer struct {
	AbstractSource   string        `json:"AbstractSource"`
	DefinitionURL    string        `json:"DefinitionURL"`
	ImageIsLogo      string        `json:"ImageIsLogo"`
	RelatedTopics    []interface{} `json:"RelatedTopics"`
	Meta             interface{}   `json:"meta"`
	DefinitionSource string        `json:"DefinitionSource"`
	Redirect         string        `json:"Redirect"`
	Infobox          string        `json:"Infobox"`
	Abstract         string        `json:"Abstract"`
	AbstractURL      string        `json:"AbstractURL"`
	Type             string        `json:"Type"`
	ImageHeight      string        `json:"ImageHeight"`
	ImageWidth       string        `json:"ImageWidth"`
	Image            string        `json:"Image"`
	AbstractText     string        `json:"AbstractText"`
	AnswerType       string        `json:"AnswerType"`
	Answer           string        `json:"Answer"`
	Heading          string        `json:"Heading"`
	Definition       string        `json:"Definition"`
	Entity           string        `json:"Entity"`
	Results          []interface{} `json:"Results"`
}

var re = regexp.MustCompile(`^Your IP address is (?P<IP>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) in .+$`)

func getIPAddress() string {
	response, err := http.Get(answerAPI)
	if err != nil {
		return ""
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}

	body := Answer{}
	json.Unmarshal(bytes, &body)

	ip := re.FindStringSubmatch(body.Answer)
	if len(ip) <= 1 {
		log.Fatal("No IP matches found for DuckDuckGo Answer")
	}

	return ip[1]
}

func main() {
	app := cli.NewApp()
	app.Name = "duckdns"
	app.Usage = "perform dynamic DNS updates in Cloudflare using DuckDuckGo Answers"
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
			EnvVar: "DUCKDNS_IP",
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
	}
	app.Run(os.Args)
}

func run(c *cli.Context) {
	key := c.String("api-key")
	email := c.String("email")
	name := c.String("name")
	domain := c.String("domain")
	ip := c.String("ip")
	if len(ip) == 0 {
		ip = getIPAddress()
	}

	fmt.Printf("Key is set to %s\n", key)
	fmt.Printf("Email is set to %s\n", email)
	fmt.Printf("IP is set to %s\n", ip)

	cf, err := cloudflare.New(key, email)
	if err != nil {
		log.Fatal("Could not initialize the Cloudflare API")
	}

	fmt.Println("Started up dynamic DNS service")

	record := cloudflare.DNSRecord{
		Type:      "A",
		Name:      name,
		Content:   ip,
		Proxiable: false,
		Proxied:   false,
		TTL:       1,
	}

	fmt.Printf("Your IP address is: %s\n", ip)

	zone, err := cf.ZoneIDByName(domain)
	if err != nil {
		log.Fatal(err)
	}

	filterRecord := cloudflare.DNSRecord{
		Name: record.Name,
	}

	records, err := cf.DNSRecords(zone, filterRecord)

	var id string
	existingRecord := cloudflare.DNSRecord{}
	for _, r := range records {
		if matchingRecord(r, record) {
			id = r.ID
			existingRecord = r
		}
	}

	if identicalRecords(record, existingRecord) {
		log.Fatal("Identical record already exists")
	}

	if len(id) > 0 {
		fmt.Println("Updating existing DNS record...")
		err = cf.UpdateDNSRecord(zone, id, record)
	} else {
		fmt.Println("Creating DNS record...")
		_, err = cf.CreateDNSRecord(zone, record)
	}

	if err != nil {
		fmt.Printf("Encountered error [%s]\n", err.Error())
	} else {
		fmt.Println("Successfully updated DNS record")
	}

}

func matchingRecord(a, b cloudflare.DNSRecord) bool {
	return false
}

func identicalRecords(a, b cloudflare.DNSRecord) bool {
	fmt.Printf("Type = %v, %v\n", a.Type, b.Type)
	fmt.Printf("Name = %v, %v\n", a.Name, b.Name)
	fmt.Printf("Content = %v, %v\n", a.Content, b.Content)
	fmt.Printf("Proxiable = %v, %v\n", a.Proxiable, b.Proxiable)
	fmt.Printf("Proxied = %v, %v\n", a.Proxied, b.Proxied)
	fmt.Printf("TTL = %v, %v\n", a.TTL, b.TTL)
	return a.Type == b.Type &&
		a.Name == b.Name &&
		a.Proxiable == b.Proxiable &&
		a.Proxied == b.Proxied &&
		a.Content == b.Content &&
		a.TTL == b.TTL
}
