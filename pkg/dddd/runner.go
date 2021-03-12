package dddd

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	ddg "github.com/tcarrio/dddd/pkg/duckduckgo"
	"github.com/tcarrio/dddd/pkg/logger"
)

// Run starts the dynamic DNS update
func Run(config Config) {
	err := ValidateConfig(&config)
	if err != nil {
		panic("Provided configuration was invalid!")
	}

	logLevel := logger.ParseLogLevel(config.LogLevel)
	logger := logger.New(logLevel)

	if len(config.IP) == 0 {
		config.IP = ddg.GetIPAddress()
	}

	if len(config.RecordType) == 0 {
		config.RecordType = "A"
	}

	logger.Debug(fmt.Sprintf("Key is set to %s", config.Key))
	logger.Debug(fmt.Sprintf("Email is set to %s", config.Email))
	logger.Debug(fmt.Sprintf("IP is set to %s", config.IP))

	cf, err := cloudflare.New(config.Key, config.Email)
	if err != nil {
		logger.Fatal("Could not initialize the Cloudflare API")
	}

	logger.Info("Started up dynamic DNS service")

	record := cloudflare.DNSRecord{
		Type:      config.RecordType,
		Name:      config.Name,
		Content:   config.IP,
		Proxiable: false,
		Proxied:   nil,
		TTL:       1,
	}

	logger.Info(fmt.Sprintf("Your IP address is: %s", config.IP))

	zone, err := cf.ZoneIDByName(config.Domain)
	if err != nil {
		logger.Warn("Failed to get zone id!")
		logger.Fatal(err)
	}

	records, err := cf.DNSRecords(context.Background(), zone, cloudflare.DNSRecord{})
	logger.Info(fmt.Sprintf("Found %d records", len(records)))

	var id string
	existingRecord := cloudflare.DNSRecord{}
	for _, r := range records {
		if MatchingRecord(r, record, config.Domain) {
			id = r.ID
			existingRecord = r
		}
	}

	if IdenticalRecords(record, existingRecord, config.Domain) {
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
