package jwtutil

import "errors"

var (
	// ErrExpired indicates that the JWT has expired.
	ErrExpired = errors.New("jwt: token is expired")
	// ErrNotYetValid indicates that the JWT is not yet valid.
	ErrNotYetValid = errors.New("jwt: token is not yet valid")
	// ErrInvalidPayload indicates that the payload is invalid.
	ErrInvalidPayload = errors.New("jwt: invalid payload")
	// ErrInvalidKey indicates that the provided secret key is invalid.
	ErrInvalidKey = errors.New("jwt: invalid secret key")
	// ErrInvalidSignature indicates that the JWT signature is invalid.
	ErrInvalidSignature = errors.New("jwt: invalid signature")
	// ErrInvalidAlgorithm indicates that the JWT algorithm is not allowed.
	ErrInvalidAlgorithm = errors.New("jwt: invalid algorithm")
)
