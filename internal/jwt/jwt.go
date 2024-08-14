package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (m *Manager) GenerateRefreshToken(guid uuid.UUID, ip string) (string, error) {
	jwtClaims := jwt.MapClaims{
		"guid": guid,
		"ip":   ip,
	}

	token := jwt.NewWithClaims(
		m.method,
		jwtClaims,
	)

	return token.SignedString([]byte(m.cfg.Key))
}

func (m *Manager) GenerateAccessToken(guid uuid.UUID, ip string, refreshId int) (string, error) {
	jwtClaims := jwt.MapClaims{
		"guid":      guid,
		"ip":        ip,
		"refreshId": refreshId,
	}

	token := jwt.NewWithClaims(
		m.method,
		jwtClaims,
	)

	return token.SignedString([]byte(m.cfg.Key))
}
