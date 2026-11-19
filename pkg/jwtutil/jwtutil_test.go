package jwtutil

import (
	"errors"
	"testing"
	"time"
)

func TestGenerateSecretKey(t *testing.T) {
	key1, err := GenerateSecretKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(key1) < 64 {
		t.Errorf("expected key length >= 64, got %d", len(key1))
	}

	key2, err := GenerateSecretKey()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if key1 == key2 {
		t.Errorf("expected different keys, got the same one: %s", key1)
	}
}

func TestGenerateJWT_Success(t *testing.T) {
	secret, _ := GenerateSecretKey()
	now := time.Now().Unix()
	payload := JwtPayload{
		Iss:  "test-issuer",
		Sub:  "user-123",
		Aud:  "test-audience",
		Exp:  now + 3600,
		Nbf:  now - 60,
		Iat:  now,
		Jti:  "some-uuid",
		Data: map[string]interface{}{"role": "admin"},
	}

	token, err := GenerateJWT(payload, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestDecodeJWT_Success(t *testing.T) {
	secret, _ := GenerateSecretKey()
	now := time.Now().Unix()
	payload := JwtPayload{
		Iss:  "test-issuer",
		Sub:  "user-123",
		Aud:  "test-audience",
		Exp:  now + 3600,
		Nbf:  now - 60,
		Iat:  now,
		Jti:  "some-uuid",
		Data: map[string]interface{}{"role": "admin"},
	}

	token, err := GenerateJWT(payload, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	decodedPayload, err := DecodeJWT(token, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if decodedPayload.Iss != payload.Iss {
		t.Errorf("expected Iss %s, got %s", payload.Iss, decodedPayload.Iss)
	}
	if decodedPayload.Sub != payload.Sub {
		t.Errorf("expected Sub %s, got %s", payload.Sub, decodedPayload.Sub)
	}
	if decodedPayload.Aud != payload.Aud {
		t.Errorf("expected Aud %s, got %s", payload.Aud, decodedPayload.Aud)
	}
	if decodedPayload.Exp != payload.Exp {
		t.Errorf("expected Exp %d, got %d", payload.Exp, decodedPayload.Exp)
	}
	if decodedPayload.Nbf != payload.Nbf {
		t.Errorf("expected Nbf %d, got %d", payload.Nbf, decodedPayload.Nbf)
	}
	if decodedPayload.Iat != payload.Iat {
		t.Errorf("expected Iat %d, got %d", payload.Iat, decodedPayload.Iat)
	}
	if decodedPayload.Jti != payload.Jti {
		t.Errorf("expected Jti %s, got %s", payload.Jti, decodedPayload.Jti)
	}
}

func TestDecodeJWT_Expired(t *testing.T) {
	secret, _ := GenerateSecretKey()
	now := time.Now().Unix()
	payload := JwtPayload{
		Exp: now - 60, // Expired
		Nbf: now - 120,
		Iat: now - 180,
		Jti: "uuid",
	}

	token, err := GenerateJWT(payload, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = DecodeJWT(token, secret)
	if !errors.Is(err, ErrExpired) {
		t.Errorf("expected ErrExpired, got %v", err)
	}
}

func TestDecodeJWT_NotYetValid(t *testing.T) {
	secret, _ := GenerateSecretKey()
	now := time.Now().Unix()
	payload := JwtPayload{
		Exp: now + 3600,
		Nbf: now + 60, // Not yet valid
		Iat: now - 60,
		Jti: "uuid",
	}

	token, err := GenerateJWT(payload, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = DecodeJWT(token, secret)
	if !errors.Is(err, ErrNotYetValid) {
		t.Errorf("expected ErrNotYetValid, got %v", err)
	}
}

func TestDecodeJWT_InvalidSignature(t *testing.T) {
	secret1, _ := GenerateSecretKey()
	secret2, _ := GenerateSecretKey()
	now := time.Now().Unix()
	payload := JwtPayload{
		Exp: now + 3600,
		Nbf: now - 60,
		Iat: now,
		Jti: "uuid",
	}

	token, err := GenerateJWT(payload, secret1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = DecodeJWT(token, secret2)
	if err == nil {
		t.Error("expected error from invalid signature, got nil")
	}
}

func TestDecodeJWT_InvalidKey(t *testing.T) {
	payload := JwtPayload{
		Exp: time.Now().Unix() + 3600,
		Jti: "uuid",
	}
	secret, _ := GenerateSecretKey()
	token, _ := GenerateJWT(payload, secret)

	_, err := DecodeJWT(token, "")
	if !errors.Is(err, ErrInvalidKey) {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}
