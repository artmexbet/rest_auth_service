package models

import "github.com/google/uuid"

type User struct {
	Guid uuid.UUID
	Ip   string
}

type AccessRefreshJSON struct {
	AccessT  string `json:"accessToken"`
	RefreshT string `json:"refreshT"`
}
