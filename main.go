package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
)

type GeoIPResponse struct {
	Country string  `json:"country"`
	ISOCode string  `json:"iso_code"`
	City    string  `json:"city"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using default environment variables")
	}

	// Get port from environment variable, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("No PORT environment variable detected, defaulting to 8080")
	}

	// Open the GeoLite2-City database
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Failed to open GeoLite2-City.mmdb file:", err)
	}
	defer db.Close()

	// Define the GeoIP handler
	http.HandleFunc("/geoip", func(w http.ResponseWriter, r *http.Request) {
		// Extract IP address from query parameters
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "Missing 'ip' query parameter", http.StatusBadRequest)
			return
		}

		// Parse the IP address
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		// Perform GeoIP lookup
		record, err := db.City(parsedIP)
		if err != nil {
			log.Printf("GeoIP lookup error for IP %s: %v", ip, err)
			http.Error(w, "Internal server error: Could not perform IP lookup", http.StatusInternalServerError)
			return
		}

		// Prepare the response
		response := GeoIPResponse{
			Country: record.Country.Names["en"],
			ISOCode: record.Country.IsoCode,
			City:    record.City.Names["en"],
			Lat:     record.Location.Latitude,
			Lon:     record.Location.Longitude,
		}

		// Handle missing city or country information
		if response.Country == "" {
			response.Country = "Unknown"
		}
		if response.City == "" {
			response.City = "Unknown"
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the HTTP server
	log.Printf("GeoIP API server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
