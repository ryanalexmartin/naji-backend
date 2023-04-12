package services

import (
	"backend/models"
	"testing"
)

type mockWebSocketConn struct {
	lastMessage string
}

func (m *mockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	m.lastMessage = string(data)
	return nil
}

func newMockWebSocketConn() *mockWebSocketConn {
	return &mockWebSocketConn{}
}

func TestMatchUsers(t *testing.T) {
	user1 := models.NewUser(nil)
	user2 := models.NewUser(nil)

	room := MatchUsers(user1, user2)

	if user1.Room != room.ID {
		t.Errorf("User1 not added to the room. Expected room ID: %s, got: %s", room.ID, user1.Room)
	}

	if user2.Room != room.ID {
		t.Errorf("User2 not added to the room. Expected room ID: %s, got: %s", room.ID, user2.Room)
	}

	if room.UserCount() != 2 {
		t.Errorf("Room does not have the correct number of users. Expected: 2, got: %d", room.UserCount())
	}
}

func TestBroadcastMessage(t *testing.T) {
	user1 := models.NewUser(nil)
	user2 := models.NewUser(nil)

	room := MatchUsers(user1, user2)

	message := "Hello, world!"

	// Mock the WebSocket connections to capture the sent messages
	user1.Conn = newMockWebSocketConn()
	user2.Conn = newMockWebSocketConn()

	BroadcastMessage(room, user1, message)

	if user1.Conn.(*mockWebSocketConn).lastMessage != message {
		t.Errorf("User1 did not receive the message. Expected: %s, got: %s", message, user1.Conn.(*mockWebSocketConn).lastMessage)
	}

	if user2.Conn.(*mockWebSocketConn).lastMessage != message {
		t.Errorf("User2 did not receive the message. Expected: %s, got: %s", message, user2.Conn.(*mockWebSocketConn).lastMessage)
	}
}

func TestLeaveChatSession(t *testing.T) {
	user1 := models.NewUser(nil)
	user2 := models.NewUser(nil)

	room := MatchUsers(user1, user2)

	LeaveChatSession(room, user1)

	if user1.Room != "" {
		t.Errorf("User1 still has a room ID after leaving. Expected: \"\", got: %s", user1.Room)
	}

	if room.UserCount() != 1 {
		t.Errorf("Room still has User1. Expected: 1, got: %d", room.UserCount())
	}

	LeaveChatSession(room, user2)

	if user2.Room != "" {
		t.Errorf("User2 still has a room ID after leaving. Expected: \"\", got: %s", user2.Room)
	}

	if room.UserCount() != 0 {
		t.Errorf("Room still has users after both left. Expected: 0, got: %d", room.UserCount())
	}
}
