package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Helper to setup handler and execute the request
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(keyHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// Test for valid request
func TestKeyHandler_ValidRequest(t *testing.T) {
	length := 8
	req, err := http.NewRequest("GET", "/key/"+strconv.Itoa(length), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := executeRequest(req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response length matches the length requested
	expectedLength := length * 2 // hex length
	if len(rr.Body.String()) != expectedLength {
		t.Errorf("handler returned unexpected body length: got %v want %v",
			len(rr.Body.String()), expectedLength)
	}
}

// Test for request exceeding max-size
func TestKeyHandler_ExceedingMaxSizeRequest(t *testing.T) {
	length := maxSize * 2
	req, err := http.NewRequest("GET", "/key/"+strconv.Itoa(length), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := executeRequest(req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Test for invalid length
func TestKeyHandler_InvalidLength(t *testing.T) {
	req, err := http.NewRequest("GET", "/key/invalid_length", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := executeRequest(req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Test for missing length
func TestKeyHandler_MissingLength(t *testing.T) {
	req, err := http.NewRequest("GET", "/key/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := executeRequest(req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
