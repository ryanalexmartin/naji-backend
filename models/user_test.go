package models

import (
	"backend/server"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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

func TestUserWebSocketConnection(t *testing.T) {
	addr := "localhost:8081"
	go server.StartServer(addr)
	time.Sleep(500 * time.Millisecond)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	header := make(http.Header)

	dialer := &websocket.Dialer{}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	user := NewUser(conn)
	if user.Conn == nil {
		t.Error("User WebSocket connection is nil")
	}

	if err := user.Conn.WriteMessage(1, []byte("Test message")); err != nil {
		t.Errorf("Failed to write message through WebSocket: %v", err)
	}
}
