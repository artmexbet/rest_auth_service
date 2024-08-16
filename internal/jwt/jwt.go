package jwt

import (
	"bytes"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"restAuthPart/internal/models"
	"time"
)

// Config ...
type Config struct {
	Key string `yaml:"key" env:"KEY" env-default:"secretkey"`
}

// Manager ...
type Manager struct {
	cfg    *Config
	method jwt.SigningMethod
}

// New ...
func New(cfg *Config, method jwt.SigningMethod) *Manager {
	return &Manager{
		cfg:    cfg,
		method: method,
	}
}

// GenerateRefreshToken generates refresh token
func (m *Manager) GenerateRefreshToken(guid uuid.UUID, ip string) (string, error) {
	jwtClaims := models.RefreshTokenClaims{
		Guid: guid,
		Ip:   ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		m.method,
		jwtClaims,
	)

	return token.SignedString([]byte(m.cfg.Key))
}

// GenerateAccessToken generates access token
func (m *Manager) GenerateAccessToken(guid uuid.UUID, ip string, id int) (string, error) {
	jwtClaims := models.AccessTokenClaims{
		Guid:      guid,
		Ip:        ip,
		RefreshId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(
		m.method,
		jwtClaims,
	)

	return token.SignedString([]byte(m.cfg.Key))
}

// GetClaims returns claims from token
func (m *Manager) GetClaims(token string, claimsType jwt.Claims) (jwt.Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claimsType, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.Key), nil
	})
	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		return parsedToken.Claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// CompareTokens compares token and hashed token
func (m *Manager) CompareTokens(token string, hashedToken []byte) bool {
	//return bcrypt.CompareHashAndPassword(hashedToken, []byte(token)) == nil
	return bytes.Compare(hashedToken, []byte(token)) == 0
}
