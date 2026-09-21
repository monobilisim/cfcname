package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	var token string
	var targetURL string
	var newURL string
	var fullURL bool

	flag.StringVar(&token, "token", "", "Cloudflare API token")
	flag.StringVar(&token, "t", "", "Cloudflare API token")

	flag.StringVar(&targetURL, "target-url", "", "Target URL to be changed")
	flag.StringVar(&targetURL, "u", "", "Target URL to be changed")

	flag.StringVar(&newURL, "new-url", "", "New URL")
	flag.StringVar(&newURL, "n", "", "New URL")

	flag.BoolVar(&fullURL, "full-url", false, "Target URL is the full URL to be replaced")
	flag.BoolVar(&fullURL, "f", false, "Target URL is the full URL to be replaced")

	flag.Parse()

	if token == "" || targetURL == "" {
		flag.Usage()
		return
	}

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	zones, err := api.ListZones(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	recordsChan := make(chan map[string][]cloudflare.DNSRecord, len(zones))

	for _, zone := range zones {
		wg.Add(1)
		go func(z cloudflare.Zone) {
			defer wg.Done()

			rc := cloudflare.ListDNSRecordsParams{
				Type: "CNAME",
			}

			records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(z.ID), rc)
			if err != nil {
				log.Printf("Error getting records for zone %s: %v", z.Name, err)
				return
			}

			recordsChan <- map[string][]cloudflare.DNSRecord{z.ID: records}
		}(zone)
	}

	go func() {
		wg.Wait()
		close(recordsChan)
	}()

	allRecords := make(map[string][]cloudflare.DNSRecord)
	containingRecords := make(map[string][]cloudflare.DNSRecord)
	for records := range recordsChan {
		for zoneID, zoneRecords := range records {
			allRecords[zoneID] = zoneRecords
		}
	}

	for zoneID, records := range allRecords {
		for _, record := range records {
			if (fullURL && record.Content == targetURL) || (!fullURL && strings.HasSuffix(record.Content, targetURL)) {
				containingRecords[zoneID] = append(containingRecords[zoneID], record)
			}
		}
	}

	for zoneID, records := range containingRecords {
		fmt.Printf("\nZone: %s\n", zoneID)
		for _, record := range records {
			if (fullURL && record.Content == targetURL) || (!fullURL && strings.HasSuffix(record.Content, targetURL)) {
				if newURL == "" {
					fmt.Printf("  Name: %s, Target: %s\n", record.Name, record.Content)
				} else {
					var nURL string
					if !fullURL {
						prefix, _ := strings.CutSuffix(record.Content, targetURL)
						nURL = prefix + newURL
					} else {
						nURL = newURL
					}
					fmt.Printf("  Before\n    Name: %s, Target: %s\n", record.Name, record.Content)

					record.Content = nURL
					newRecord, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
						ID:      record.ID,
						Content: record.Content,
					})
					if err != nil {
						fmt.Printf("Error updating record %s: %v\n", record.Name, err)
					}
					fmt.Printf("  After\n    Name: %s, Target: %s\n", newRecord.Name, newRecord.Content)
				}
			}
		}
	}
}
