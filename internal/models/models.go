package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	Guid  uuid.UUID
	Ip    string
	Email string
}

type AccessRefreshJSON struct {
	AccessT  string `json:"accessT"`
	RefreshT string `json:"refreshT"`
}

type RefreshTokenJSON struct {
	RefreshT string `json:"refreshT"`
	AccessT  string `json:"accessT"`
}

type RefreshTokenClaims struct {
	Guid uuid.UUID `json:"guid"`
	Ip   string    `json:"ip"`
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	Guid      uuid.UUID `json:"guid"`
	Ip        string    `json:"ip"`
	RefreshId int       `json:"refreshId"`
	jwt.RegisteredClaims
}
