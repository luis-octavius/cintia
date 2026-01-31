package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT_Success(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Ponga la Ponga"
	expiresIn := 24 * time.Hour

	token, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT Error: %v", err)
	}

	validatedUUID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT Error: %v", err)
	}

	if validatedUUID != id {
		t.Errorf("expected id %v to match %v", id, validatedUUID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Alright"
	expiresIn := time.Duration(0)

	token, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT Error: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected error on validating token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Sun is rising"
	expiresIn := 2 * time.Hour

	token, err := MakeJWT(id, "Sun is dying", expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT Error: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected error on validating token with wrong secret, got nil")
	}
}
