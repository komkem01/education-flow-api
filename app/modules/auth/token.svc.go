package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type accessClaims struct {
	Sub      string `json:"sub"`
	Role     string `json:"role"`
	SchoolID string `json:"school_id"`
	JTI      string `json:"jti"`
	IAT      int64  `json:"iat"`
	EXP      int64  `json:"exp"`
}

func (s *Service) createAccessToken(claims *accessClaims) (string, error) {
	headerJSON, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	unsigned := header + "." + payload

	sig := signHMAC(unsigned, s.conf.JWTSecret)
	return unsigned + "." + sig, nil
}

func (s *Service) parseAccessToken(token string) (*accessClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSig := signHMAC(unsigned, s.conf.JWTSecret)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	claims := new(accessClaims)
	if err := json.Unmarshal(payload, claims); err != nil {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}
	if claims.EXP <= time.Now().Unix() {
		return nil, fmt.Errorf("%w", ErrAuthUnauthorized)
	}

	return claims, nil
}

func signHMAC(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func newRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:])
}
