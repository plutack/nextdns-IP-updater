package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url := os.Getenv("NEXTDNS_IP")
	if url == "" {
		log.Fatal("NEXTDNS_IP environment variable is not set")
	}

	log.Printf("Starting NextDNS IP updater for: %s", url)

	intervalStr := os.Getenv("UPDATE_INTERVAL")
	interval := 10 * time.Minute
	if intervalStr != "" {
		if parsed, err := time.ParseDuration(intervalStr); err == nil {
			interval = parsed
		} else {
			log.Printf("Invalid UPDATE_INTERVAL format, using default 10m: %v", err)
		}
	}
	log.Printf("Update interval set to: %v", interval)

	for {
		updateIP(url)
		time.Sleep(interval)
	}
}

func updateIP(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error updating IP: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to update IP. Status: %s, Body: %s", resp.Status, string(body))
	} else {
		log.Printf("IP updated successfully. Response: %s", string(body))
	}
}
