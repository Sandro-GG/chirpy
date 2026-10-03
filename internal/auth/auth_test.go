package auth

import (
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "ziGGyPlayedGuitar"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		t.Fatalf("Failed to compare password: %v", err)
	}
	if !match {
		t.Error("The password does not match the generated hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "ziGGyPlayedGuitar"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Fatalf("Failed setup: could not create hash: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash errored unexpectedly: %v", err)
	}
	if !match {
		t.Errorf("Expected correct password to match, but it failed")
	}

	wrongMatch, err := CheckPasswordHash("wrong-password-223", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash errored unexpectedly: %v", err)
	}
	if wrongMatch {
		t.Errorf("Expected wrong password to fail, but it matched")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "very-ver-y-sec-ret-key-223"
	duration := 1 * time.Hour

	tokenStr, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT returned an unexpected error: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse the generated token: %v", err)
	}

	if claims.Issuer != "chirpy-access" {
		t.Errorf("Expected issuer 'chirpy-access', got '%s'", claims.Issuer)
	}

	if claims.Subject != userID.String() {
		t.Errorf("Expected subject'%s', got '%s'", userID.String(), claims.Subject)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		t.Errorf("Token is already expired")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "big-ben-bong-bang-boeing"

	tokenStr, err := MakeJWT(userID, secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make token")
	}

	id, err := ValidateJWT(tokenStr, secret)
	if err != nil {
		t.Fatalf("Failed to validate token")
	}
	if id != userID {
		t.Errorf("Expected %v, got %v", userID, id)
	}

	tokenStr, err = MakeJWT(userID, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to make token")
	}

	_, err = ValidateJWT(tokenStr, secret)
	if err == nil {
		t.Errorf("Expected an error validating an expired token")
	}

	tokenStr, err = MakeJWT(userID, "diff-ere-n-t-secret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to make token")
	}

	_, err = ValidateJWT(tokenStr, secret)
	if err == nil {
		t.Errorf("Expected an error validating a token signed with the wrong secret")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := map[string][]string{
		"Testing":       {"hello hi", "boom boom"},
		"Spatula":       {"big", "bob", "squid clarinet"},
		"Authorization": {"Bearer fasdck-fsdakv-fkadqew-fdass", "mod cod pod"},
	}

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Failed to get token")
	}

	if token != "fasdck-fsdakv-fkadqew-fdass" {
		t.Errorf("Failed to extract token correctly")
	}
}
