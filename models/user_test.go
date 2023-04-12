package models

import (
	"github.com/gorilla/websocket"
	"testing"
)

func TestUserCreation(t *testing.T) {
	// Dummy WebSocket connection
	conn := &websocket.Conn{}

	// Create a new user
	user := NewUser(conn)

	// Check if the user has a unique ID
	if user.ID == "" {
		t.Error("User ID should not be empty")
	}

	// Check if the user's WebSocket connection is valid
	if user.Conn != conn {
		t.Error("User connection is not set correctly")
	}

	// Check if the user's Room field is initially empty
	if user.Room != "" {
		t.Error("User's room ID should be empty initially")
	}
}
