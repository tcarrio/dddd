package dddd

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	ddg "github.com/tcarrio/dddd/pkg/duckduckgo"
	"github.com/tcarrio/dddd/pkg/logger"
)

// RunnerArgs are passed to run the main functionality of dddd
type RunnerArgs struct {
	Key        string
	Email      string
	Name       string
	Domain     string
	IP         string
	RecordType string
	LogLevel   string
}

// Run starts the dynamic DNS update
func Run(args RunnerArgs) {
	logLevel := logger.ParseLogLevel(args.LogLevel)
	logger := logger.New(logLevel)

	if len(args.IP) == 0 {
		args.IP = ddg.GetIPAddress()
	}

	if len(args.RecordType) == 0 {
		args.RecordType = "A"
	}

	logger.Debug(fmt.Sprintf("Key is set to %s", args.Key))
	logger.Debug(fmt.Sprintf("Email is set to %s", args.Email))
	logger.Debug(fmt.Sprintf("IP is set to %s", args.IP))

	cf, err := cloudflare.New(args.Key, args.Email)
	if err != nil {
		logger.Fatal("Could not initialize the Cloudflare API")
	}

	logger.Info("Started up dynamic DNS service")

	record := cloudflare.DNSRecord{
		Type:      args.RecordType,
		Name:      args.Name,
		Content:   args.IP,
		Proxiable: false,
		Proxied:   nil,
		TTL:       1,
	}

	logger.Info(fmt.Sprintf("Your IP address is: %s", args.IP))

	zone, err := cf.ZoneIDByName(args.Domain)
	if err != nil {
		logger.Warn("Failed to get zone id!")
		logger.Fatal(err)
	}

	records, err := cf.DNSRecords(context.Background(), zone, cloudflare.DNSRecord{})
	logger.Info(fmt.Sprintf("Found %d records", len(records)))

	var id string
	existingRecord := cloudflare.DNSRecord{}
	for _, r := range records {
		if MatchingRecord(r, record, args.Domain) {
			id = r.ID
			existingRecord = r
		}
	}

	if IdenticalRecords(record, existingRecord, args.Domain) {
		logger.Fatal("Identical record already exists")
	}

	if len(id) > 0 {
		logger.Info(fmt.Sprintf("Updating existing DNS record %s...", id))
		err = cf.UpdateDNSRecord(context.Background(), zone, id, record)
	} else {
		logger.Info("Creating DNS record...")
		_, err = cf.CreateDNSRecord(context.Background(), zone, record)
	}

	if err != nil {
		logger.Info(fmt.Sprintf("Encountered error [%s]", err.Error()))
	} else {
		logger.Info("Successfully updated DNS record")
	}
}
