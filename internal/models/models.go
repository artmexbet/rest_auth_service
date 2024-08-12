package models

import "github.com/google/uuid"

type User struct {
	Guid uuid.UUID
	Ip   string
}
