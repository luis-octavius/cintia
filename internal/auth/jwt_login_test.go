package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidateJWT_Success(t *testing.T) {
	id := uuid.New()
	secret := "Ponga la Ponga"
	email := "test@example.com"
	role := "admin"

	token, err := MakeJWT(id, email, role, secret)
	if err != nil {
		t.Fatalf("MakeJWT Error: %v", err)
	}

	claims, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT Error: %v", err)
	}

	if claims.UserID != id.String() {
		t.Errorf("expected id %v to match %v", id.String(), claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected email %v to match %v", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("expected role %v to match %v", role, claims.Role)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	id := uuid.New()
	secret := "Alright"
	email := "test@example.com"
	role := "user"

	claims := UserClaims{
		UserID: id.String(),
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "cintia",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Token signing error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("Expected error on validating token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	id := uuid.New()
	secret := "Sun is rising"
	email := "test@example.com"
	role := "user"

	token, err := MakeJWT(id, email, role, "Sun is dying")
	if err != nil {
		t.Fatalf("MakeJWT Error: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("Expected error on validating token with wrong secret, got nil")
	}
}
