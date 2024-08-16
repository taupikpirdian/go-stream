package entity

import "github.com/google/uuid"

type User struct {
	ID                uuid.UUID
	Name              string
	PricePlan         string
	Email             string
	IsEligibelChatBot bool
}
