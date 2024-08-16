package entity

import "github.com/google/uuid"

type WhitelistChat struct {
	ID    uuid.UUID
	Email string
}
