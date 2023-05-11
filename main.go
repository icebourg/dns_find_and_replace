package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

type DeleteRecord struct {
	recordId string
	zoneId   string
	record   string
}

func main() {
	findValue := flag.String("find", "REQUIRED", "the DNS record content we are searching for")
	replacementValue := flag.String("replacement", "REQUIRED", "the content of the record we want to replace with")
	replacementType := flag.String("replacementtype", "CNAME", "the type of the record we want to create")
	batchSize := flag.Int("batchsize", 100, "the number of DNS records we will operate on at a given time")
	flag.Parse()

	if *findValue == "REQUIRED" || *replacementValue == "REQUIRED" || *replacementType == "REQUIRED" {
		log.Fatal("Please pass required find and replacement values")
	}

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch Zones
	zones, err := api.ListZones(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var recordsToDelete []DeleteRecord

	for _, z := range zones {
		records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(z.ID), cloudflare.ListDNSRecordsParams{Content: *findValue})
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range records {
			fmt.Printf("Found record %s -> %s (%s)\n", r.Name, r.Content, r.Type)
			recordsToDelete = append(recordsToDelete, DeleteRecord{recordId: r.ID, zoneId: z.ID, record: r.Name})

			if len(recordsToDelete) >= *batchSize {
				replaceRecords(ctx, api, recordsToDelete, *replacementValue, *replacementType)
				recordsToDelete = nil
			}
		}
	}

	// anymore records to delete? We've gone through all zones by now, so these are the last of it
	if len(recordsToDelete) > 0 {
		replaceRecords(ctx, api, recordsToDelete, *replacementValue, *replacementType)
		recordsToDelete = nil
	}
}

func replaceRecords(ctx context.Context, api *cloudflare.API, recordsToDelete []DeleteRecord, replacementValue string, replacementType string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Ready to commit %d replacements? Yes to proceed, any other value to stop...", len(recordsToDelete))
	s, _ := r.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(s)) != "yes" {
		os.Exit(5)
	}

	for _, record := range recordsToDelete {
		fmt.Printf("Replacing %s\n", record.record)
		err := api.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(record.zoneId), record.recordId)
		if err != nil {
			log.Fatal(err)
		}

		_, err2 := api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(record.zoneId), cloudflare.CreateDNSRecordParams{Type: replacementType, Name: record.record, Content: replacementValue, ZoneID: record.zoneId})
		if err2 != nil {
			fmt.Printf("WARNING! This record has been deleted but I failed to create the replacement record! Please handle manually for record %s\n\n", record.record)
			log.Fatal(err)
		}
	}
}
