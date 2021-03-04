package dddd

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

// MatchingRecord checks that two Cloudflare records have the same uniqueness signature
func MatchingRecord(a, b cloudflare.DNSRecord, domain string) bool {
	return a.Type == b.Type && (a.Name == b.Name || a.Name == fmt.Sprintf("%s.%s", b.Name, domain))
}

// IdenticalRecords checks that two Cloudflare records have entirely identical configurations
func IdenticalRecords(a, b cloudflare.DNSRecord, domain string) bool {
	return MatchingRecord(a, b, domain) &&
		a.Proxiable == b.Proxiable &&
		a.Proxied == b.Proxied &&
		a.Content == b.Content &&
		a.TTL == b.TTL
}
