package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	loadConfig()

	url := os.Getenv("NEXTDNS_IP")
	if url == "" {
		log.Fatal("NEXTDNS_IP environment variable is not set")
	}

	log.Printf("Starting NextDNS IP updater for: %s", url)

	intervalStr := os.Getenv("UPDATE_INTERVAL")
	if intervalStr == "" {
		intervalStr = "10m"
	}
	interval := 10 * time.Minute
	if parsed, err := time.ParseDuration(intervalStr); err == nil {
		interval = parsed
	} else {
		log.Printf("Invalid UPDATE_INTERVAL format, using default 10m: %v", err)
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

func loadConfig() {
	// 1. Check user's XDG config path (~/.config/nextdns-ip/config)
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		if usr, err := user.Current(); err == nil {
			xdgConfigHome = filepath.Join(usr.HomeDir, ".config")
		}
	}

	configPaths := []string{}
	if xdgConfigHome != "" {
		configPaths = append(configPaths, filepath.Join(xdgConfigHome, "nextdns-ip", "config"))
	}

	// 2. Check system-wide config path (/etc/nextdns-ip.conf)
	configPaths = append(configPaths, "/etc/nextdns-ip.conf")

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("Loading configuration from: %s", path)
			parseConfigFile(path)
			// If we found a config file, stop searching (user config overrides system)
			return
		}
	}
}

func parseConfigFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Warning: Failed to open config file %s: %v", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if len(value) >= 2 && (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") || strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}
