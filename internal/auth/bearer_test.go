package auth

import (
	"net/http"
	"testing"
)

var Client = &http.Client{}

func TestGetBearerTokenSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "http://www.example.com", nil)

	expected := "okdokie"
	req.Header.Set("Authorization", "Bearer okdokie")

	bearer, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("Error in GetBearerToken: %v", err)
	}

	if expected != bearer {
		t.Errorf("expected %v, got %v", expected, bearer)
	}

}
