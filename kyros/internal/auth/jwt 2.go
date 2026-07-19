package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWKS is the JSON Web Key Set shape Keycloak returns from its certs endpoint.
type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Validator validates Keycloak-issued RS256 access tokens against the realm's
// JWKS endpoint. It is safe for concurrent use.
//
// Usage:
//
//	v := auth.NewValidator("https://sso.example.com/realms/kyros")
//	claims, err := v.ValidateToken(bearerString)
type Validator struct {
	issuerURL  string
	jwksURL    string
	httpClient *http.Client
	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	lastFetch  time.Time
}

func NewValidator(issuerURL string) *Validator {
	return &Validator{
		issuerURL:  issuerURL,
		jwksURL:    fmt.Sprintf("%s/protocol/openid-connect/certs", issuerURL),
		httpClient: &http.Client{Timeout: 5 * time.Second},
		keys:       make(map[string]*rsa.PublicKey),
	}
}

func (v *Validator) fetchJWKS() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Rate limit fetching to once per minute.
	if time.Since(v.lastFetch) < time.Minute {
		return nil
	}

	resp, err := v.httpClient.Get(v.jwksURL)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			continue
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
		newKeys[key.Kid] = pubKey
	}

	v.keys = newKeys
	v.lastFetch = time.Now()
	return nil
}

func (v *Validator) getKey(token *jwt.Token) (interface{}, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	v.mu.RLock()
	key, exists := v.keys[kid]
	v.mu.RUnlock()

	if !exists {
		// Try fetching new keys in case the realm rotated.
		if err := v.fetchJWKS(); err != nil {
			return nil, err
		}

		v.mu.RLock()
		key, exists = v.keys[kid]
		v.mu.RUnlock()

		if !exists {
			return nil, fmt.Errorf("key %s not found", kid)
		}
	}

	return key, nil
}

// ValidateToken parses, signature-verifies, issuer-validates, and expiry-checks
// the supplied bearer token. On success it returns the claims as a map so
// callers don't need to depend on golang-jwt types.
func (v *Validator) ValidateToken(tokenString string) (map[string]any, error) {
	if v.issuerURL == "" {
		return nil, errors.New("validator has no issuer configured")
	}

	// First fetch if we have no keys.
	v.mu.RLock()
	keysCount := len(v.keys)
	v.mu.RUnlock()

	if keysCount == 0 {
		if err := v.fetchJWKS(); err != nil {
			return nil, fmt.Errorf("initial JWKS fetch failed: %w", err)
		}
	}

	token, err := jwt.Parse(tokenString, v.getKey, jwt.WithValidMethods([]string{"RS256"}))
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Issuer validation
	if iss, ok := claims["iss"].(string); !ok || iss != v.issuerURL {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}
