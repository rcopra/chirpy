package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == "correct horse battery staple" {
		t.Fatal("HashPassword returned the plaintext password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "correct horse battery staple"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	tests := []struct {
		name      string
		password  string
		hash      string
		wantMatch bool
		wantErr   bool
	}{
		{
			name:      "correct password",
			password:  password,
			hash:      hash,
			wantMatch: true,
		},
		{
			name:     "wrong password",
			password: "wrong password",
			hash:     hash,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "not-a-valid-hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckPasswordHash error = %v, wantErr %v", err, tt.wantErr)
			}
			if match != tt.wantMatch {
				t.Errorf("CheckPasswordHash match = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tests := []struct {
		name           string
		makeToken      func(t *testing.T) string
		validateSecret string
		wantUserID     uuid.UUID
		wantErr        bool
	}{
		{
			name: "valid token",
			makeToken: func(t *testing.T) string {
				return mustMakeJWT(t, userID, secret, time.Hour)
			},
			validateSecret: secret,
			wantUserID:     userID,
		},
		{
			name: "wrong secret",
			makeToken: func(t *testing.T) string {
				return mustMakeJWT(t, userID, secret, time.Hour)
			},
			validateSecret: "wrong-secret",
			wantErr:        true,
		},
		{
			name: "expired token",
			makeToken: func(t *testing.T) string {
				return mustMakeJWT(t, userID, secret, -time.Hour)
			},
			validateSecret: secret,
			wantErr:        true,
		},
		{
			name: "malformed token",
			makeToken: func(t *testing.T) string {
				return "not.a.jwt"
			},
			validateSecret: secret,
			wantErr:        true,
		},
		{
			name: "non-uuid subject",
			makeToken: func(t *testing.T) string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
					Subject:   "not-a-uuid",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				})
				signed, err := token.SignedString([]byte(secret))
				if err != nil {
					t.Fatalf("signing token: %v", err)
				}
				return signed
			},
			validateSecret: secret,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.makeToken(t), tt.validateSecret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if gotUserID != uuid.Nil {
					t.Errorf("ValidateJWT userID = %v, want uuid.Nil on error", gotUserID)
				}
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT userID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func mustMakeJWT(t *testing.T, userID uuid.UUID, secret string, expiresIn time.Duration) string {
	t.Helper()
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	return token
}
