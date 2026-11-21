package jwtutil

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JwtPayload represents the JWT claims.
type JwtPayload struct {
	// Iss represents the issuer.
	Iss string `json:"iss,omitempty"`

	// Sub represents the subject.
	Sub string `json:"sub,omitempty"`

	// Aud represents the audience.
	Aud string `json:"aud,omitempty"`

	// Exp represents the expiration time (Unix Timestamp in seconds).
	Exp int64 `json:"exp"`

	// Nbf represents the not before time (Unix Timestamp in seconds).
	Nbf int64 `json:"nbf"`

	// Iat represents the issued at time (Unix Timestamp in seconds).
	Iat int64 `json:"iat"`

	// Jti represents the JWT ID (UUID v4).
	Jti string `json:"jti"`

	// Data represents application-specific data.
	Data any `json:"data,omitempty"`
}

// Validate implements the jwt.Claims interface.
func (p JwtPayload) Validate() error {
	now := time.Now().Unix()

	if p.Exp != 0 && now > p.Exp {
		return ErrExpired
	}
	if p.Nbf != 0 && now < p.Nbf {
		return ErrNotYetValid
	}
	return nil
}

func (p JwtPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	if p.Exp == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(p.Exp, 0)), nil
}

func (p JwtPayload) GetNotBefore() (*jwt.NumericDate, error) {
	if p.Nbf == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(p.Nbf, 0)), nil
}

func (p JwtPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	if p.Iat == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(p.Iat, 0)), nil
}

func (p JwtPayload) GetIssuer() (string, error) {
	return p.Iss, nil
}

func (p JwtPayload) GetSubject() (string, error) {
	return p.Sub, nil
}

func (p JwtPayload) GetAudience() (jwt.ClaimStrings, error) {
	if p.Aud == "" {
		return nil, nil
	}
	return jwt.ClaimStrings{p.Aud}, nil
}

// GenerateSecretKey generates a cryptographically secure random secret key as a hex string.
func GenerateSecretKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateJWT generates a new JWT token with the given payload and secret.
func GenerateJWT(payload JwtPayload, secret string) (string, error) {
	if secret == "" {
		return "", ErrInvalidKey
	}

	if payload.Jti == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("failed to generate UUID: %w", err)
		}
		payload.Jti = id.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// DecodeJWT decodes and validates a JWT token.
func DecodeJWT(tokenString string, secret string) (*JwtPayload, error) {
	if secret == "" {
		return nil, ErrInvalidKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidAlgorithm, token.Header["alg"])
		}
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidAlgorithm
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JwtPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidPayload
}
