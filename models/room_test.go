package models

import (
	"testing"
)

func TestRoomCreation(t *testing.T) {
	// Create a new room
	room := NewRoom()

	// Check if the room has a unique ID
	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}

	// Count the number of users in the room
	userCount := 0
	room.Users.Range(func(_, _ interface{}) bool {
		userCount++
		return true
	})

	// Check if the room's Users map is initially empty
	if userCount != 0 {
		t.Error("Room should have no users initially")
	}
}
