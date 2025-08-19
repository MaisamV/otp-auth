package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/otp-auth/internal/application/ports/services"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/otp-auth/pkg/errors"
)

// ECDSAJWTService implements JWTService using ECDSA algorithm
type ECDSAJWTService struct {
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

// JWTConfig holds configuration for JWT service
type JWTConfig struct {
	PrivateKeyPEM   string
	PublicKeyPEM    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour, // 7 days
		Issuer:          "otp-auth-service",
	}
}

// NewECDSAJWTService creates a new ECDSA JWT service
func NewECDSAJWTService(config JWTConfig) (*ECDSAJWTService, error) {
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey
	var err error

	if config.PrivateKeyPEM != "" && config.PublicKeyPEM != "" {
		// Parse provided keys
		privateKey, err = parsePrivateKeyFromPEM(config.PrivateKeyPEM)
		if err != nil {
			return nil, errors.NewInternalError("Failed to parse private key", err)
		}

		publicKey, err = parsePublicKeyFromPEM(config.PublicKeyPEM)
		if err != nil {
			return nil, errors.NewInternalError("Failed to parse public key", err)
		}
	} else {
		// Generate new key pair
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, errors.NewInternalError("Failed to generate ECDSA key pair", err)
		}
		publicKey = &privateKey.PublicKey
	}

	return &ECDSAJWTService{
		privateKey:      privateKey,
		publicKey:       publicKey,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.Issuer,
	}, nil
}

// CustomClaims wraps JWTClaims for jwt library compatibility
type CustomClaims struct {
	// Standard JWT claims with explicit JSON tags
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	ID        string `json:"jti"`
	// Custom claims
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
}

// Implement jwt.Claims interface methods
func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if c.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if c.IssuedAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c CustomClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c CustomClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// GenerateToken generates a JWT token from claims
func (j *ECDSAJWTService) GenerateToken(claims *services.JWTClaims) (string, error) {
	// Debug logging - function entry
	fmt.Printf("[DEBUG] GenerateToken called\n")
	fmt.Println("[DEBUG] JWT Claims - Subject:", claims.Subject, "Issuer:", claims.Issuer, "IssuedAt:", claims.IssuedAt, "ExpiresAt:", claims.ExpiresAt, "TokenID:", claims.TokenID, "ClientID:", claims.ClientID, "Scopes:", claims.Scopes)

	// Create custom claims for jwt library
	customClaims := &CustomClaims{
		Subject:   claims.Subject,
		Issuer:    claims.Issuer,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
		ID:        claims.TokenID,
		ClientID:  claims.ClientID,
		Scopes:    claims.Scopes,
	}

	// Debug logging for custom claims
	fmt.Printf("[DEBUG] Custom Claims - Subject: %s, Issuer: %s, IssuedAt: %d, ExpiresAt: %d, ID: %s, ClientID: %s, Scopes: %v\n",
		customClaims.Subject, customClaims.Issuer, customClaims.IssuedAt, customClaims.ExpiresAt, customClaims.ID, customClaims.ClientID, customClaims.Scopes)

	token := jwt.NewWithClaims(jwt.SigningMethodES256, customClaims)
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", errors.NewInternalError("Failed to sign JWT token", err)
	}

	return tokenString, nil
}

// VerifyToken verifies and parses a JWT token, returning the claims
func (j *ECDSAJWTService) VerifyToken(tokenString string) (*services.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid token", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Convert CustomClaims back to JWTClaims
		jwtClaims := &services.JWTClaims{
			Subject:   claims.Subject,
			Issuer:    claims.Issuer,
			IssuedAt:  claims.IssuedAt,
			ExpiresAt: claims.ExpiresAt,
			TokenID:   claims.ID,
			ClientID:  claims.ClientID,
			Scopes:    claims.Scopes,
		}
		return jwtClaims, nil
	}

	return nil, errors.NewUnauthorizedError("Invalid token claims", nil)
}

// GetPublicKeyPEM returns the public key in PEM format
func (j *ECDSAJWTService) GetPublicKeyPEM() (string, error) {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(j.publicKey)
	if err != nil {
		return "", errors.NewInternalError("Failed to marshal public key", err)
	}

	pemEncodedPub := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509EncodedPub,
	})

	return string(pemEncodedPub), nil
}

// GetPrivateKeyPEM returns the private key in PEM format
func (j *ECDSAJWTService) GetPrivateKeyPEM() (string, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(j.privateKey)
	if err != nil {
		return "", errors.NewInternalError("Failed to marshal private key", err)
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: x509Encoded,
	})

	return string(pemEncoded), nil
}

// Helper functions for key parsing
func parsePrivateKeyFromPEM(pemStr string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func parsePublicKeyFromPEM(pemStr string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return publicKey, nil
}
