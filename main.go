package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPResponse struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
}

func main() {
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
			City:    record.City.Names["en"],
			Lat:     record.Location.Latitude,
			Lon:     record.Location.Longitude,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("GeoIP API server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
