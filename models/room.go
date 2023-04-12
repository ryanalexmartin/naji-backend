package models

import (
	"github.com/google/uuid"
	"sync"
)

type Room struct {
	ID    string
	Users sync.Map // Use sync.Map for concurrent-safe user management
}

func NewRoom() *Room {
	return &Room{
		ID: uuid.New().String(),
	}
}
