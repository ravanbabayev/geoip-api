package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoIPHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/geoip?ip=8.8.8.8", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main()
	})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}
}
