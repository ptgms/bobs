package main

import (
	"bobs/helpers"
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/go-co-op/gocron"
	probing "github.com/prometheus-community/pro-bing"
	"sort"
	"time"
)

var triggered = false

func task() {
	sort.Slice(config.Servers[:], func(i, j int) bool {
		return config.Servers[i].Priority < config.Servers[j].Priority
	})
	//	we go down the list of servers and return if we find one that is up
	for _, server := range config.Servers {
		if !triggered {
			pinger, err := probing.NewPinger(server.Host)
			helpers.HandleError(err, false)
			pinger.SetPrivileged(true)
			pinger.Timeout = 1 * time.Second
			pinger.Count = 1
			pinger.Interval = 1 * time.Second
			err = pinger.Run()
			if err != nil {
				println(err.Error())
			} else {
				return
			}
		} else {
			//	if we get here, we have already triggered the DNS update, and we just need to check if the server is
			//	back up, so we can revert the DNS records
			pinger, err := probing.NewPinger(server.Host)
			helpers.HandleError(err, false)
			pinger.SetPrivileged(true)
			pinger.Timeout = 1 * time.Second
			pinger.Count = 1
			pinger.Interval = 1 * time.Second
			err = pinger.Run()
			if err != nil {
				println(err.Error())
			} else {
				println("Server is back up! Reverting DNS records.")
				for _, domain := range config.Domains {
					setIP(domain.Name, domain.Zone, server.Host)
				}
				triggered = false
				return
			}
		}
	}
	println("All servers are down! Updating DNS records.")
	//	if we get here, all servers are down, and we need to update the DNS records to point to the server
	//	running this program
	for _, domain := range config.Domains {
		setIP(domain.Name, domain.Zone, config.IP)
	}
}

func prepareScheduler() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Minutes().Do(task)
	if err != nil {
		return
	}
	if err != nil {
		println(err.Error())
		return
	}

	s.StartBlocking()
}

func setIP(domain string, zone string, newIP string) {
	zoneID, err := api.ZoneIDByName(zone)
	helpers.HandleError(err, false)
	records, _, _ := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})

	for _, record := range records {
		if (record.Type != "A" && record.Type != "AAAA") || record.Name != domain {
			continue
		}
		dnsRecord, err := api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Content: newIP,
		})
		helpers.HandleError(err, false)
		println(dnsRecord.Content)
	}
}

var config helpers.Config

var api *cloudflare.API

func main() {
	println("Starting BOBS!")
	config = helpers.LoadConfig()
	api, _ = cloudflare.NewWithAPIToken(config.CloudflareAPIToken)
	prepareScheduler()
}
