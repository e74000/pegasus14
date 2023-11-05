package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

var secretKey = []byte("yourSecretKey")

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claim struct {
	Email      string    `json:"email"`
	ValidUntil time.Time `json:"valid_until"`
	Signature  string    `json:"signature"`
}

func SignClaim(email string, validUntil time.Time) (string, error) {
	claim := Claim{
		Email:      email,
		ValidUntil: validUntil,
	}

	// Convert the claim to JSON
	claimJSON, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	// Create an HMAC with SHA-256
	h := hmac.New(sha256.New, secretKey)
	h.Write(claimJSON)

	// Get the HMAC signature
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Update the claim with the signature
	claim.Signature = signature

	// Convert the updated claim to JSON
	signedClaimJSON, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	return string(signedClaimJSON), nil
}

func VerifyClaim(claim Claim) (bool, error) {
	// Extract the provided signature
	providedSignature := claim.Signature
	claim.Signature = "" // Clear the signature for signature verification

	// Convert the claim to JSON
	claimJSON, err := json.Marshal(claim)
	if err != nil {
		return false, err
	}

	// Create an HMAC with SHA-256
	h := hmac.New(sha256.New, secretKey)
	h.Write(claimJSON)

	// Get the HMAC signature
	calculatedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Compare the provided and calculated signatures
	if providedSignature != calculatedSignature {
		return false, errors.New("signature verification failed")
	}

	// Verify if the claim is still valid
	currentTime := time.Now()
	if claim.ValidUntil.Before(currentTime) {
		return false, errors.New("claim has expired")
	}

	return true, nil
}
