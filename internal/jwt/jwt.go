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

func (m *Manager) GenerateToken(guid uuid.UUID, ip string) (string, error) {
	jwtClaims := jwt.MapClaims{
		"guid": guid,
		"ip":   ip,
	}

	token := jwt.NewWithClaims(
		m.method,
		jwtClaims,
	)
	// TODO: Закодировать внутрь JWT Access идентификатор refresh токена в бд, чтобы только с ним можно было рефрешить

	return token.SignedString(m.cfg.Key)
}
