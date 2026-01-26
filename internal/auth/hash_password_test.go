package auth

import (
	"testing"
)

func TestPasswordHashingSuccess(t *testing.T) {
	password := "cradle filth"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !match {
		t.Errorf("expected match to be %v, got %v", true, match)
	}
}

func TestPasswordHashingFail(t *testing.T) {
	password := "parabola"

	hashedPassword, err := HashPassword("parabol")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if match {
		t.Errorf("expected match to be %v, got %v", false, match)
	}
}
