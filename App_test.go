package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestBot(t *testing.T) {
	// Create a request
	r, err := http.NewRequest("GET", "http://test.com", nil)

	if err != nil {
		t.Fatal(err.Error())
	}
	// Handle request and store result in w
	w := httptest.NewRecorder()

	requestHandler(w, r)

	if w.Code != http.StatusOK {
		t.Fatal(w.Code, w.Body.String())
	}
}
