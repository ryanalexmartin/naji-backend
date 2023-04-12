package models

import (
	"github.com/gorilla/websocket"
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

func TestAddUserToRoom(t *testing.T) {
	// Create a new room
	room := NewRoom()

	// Create a new user with a dummy WebSocket connection
	conn := &websocket.Conn{}
	user := NewUser(conn)

	// Add the user to the room
	room.AddUser(user)

	// Check if the user is in the room
	if _, ok := room.Users.Load(user.ID); !ok {
		t.Error("User should be added to the room")
	}

	// Check if the user's room ID is updated
	if user.Room != room.ID {
		t.Error("User's room ID should be updated after adding to the room")
	}
}

func TestRemoveUserFromRoom(t *testing.T) {
	// Create a new room
	room := NewRoom()

	// Create a new user with a dummy WebSocket connection
	conn := &websocket.Conn{}
	user := NewUser(conn)

	// Add the user to the room
	room.AddUser(user)

	// Remove the user from the room
	room.RemoveUser(user.ID)

	// Check if the user is removed from the room
	if _, ok := room.Users.Load(user.ID); ok {
		t.Error("User should be removed from the room")
	}

	// Check if the user's room ID is updated
	if user.Room != "" {
		t.Error("User's room ID should be empty after removing from the room")
	}
}
