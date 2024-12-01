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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Failed to open MMDB file:", err)
	}
	defer db.Close()

	http.HandleFunc("/geoip", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "Missing 'ip' query parameter", http.StatusBadRequest)
			return
		}

		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		record, err := db.City(parsedIP)
		if err != nil {
			http.Error(w, "Failed to lookup IP address", http.StatusInternalServerError)
			return
		}

		response := GeoIPResponse{
			Country: record.Country.Names["en"],
			ISOCode: record.Country.IsoCode,
			City:    record.City.Names["en"],
			Lat:     record.Location.Latitude,
			Lon:     record.Location.Longitude,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Printf("GeoIP API server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
