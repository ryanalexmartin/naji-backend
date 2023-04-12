package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketUpgrade(t *testing.T) {
	// Create an HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(upgradeToWebSocket))
	defer ts.Close()

	// Prepare a WebSocket connection request
	u := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Make a request to the test server
	_, resp, err := websocket.DefaultDialer.Dial(u, nil)

	// Check if there's no error while making a request
	if err != nil && err != websocket.ErrBadHandshake {
		t.Fatalf("Error making a request to the test server: %v", err)
	}

	// Check if the response status code is 101 (Switching Protocols)
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Expected status 101, got %d", resp.StatusCode)
	}
}

func TestHandleWebSocketMessage(t *testing.T) {
	// Create an HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		handleWebSocketMessage(conn)
	}))
	defer ts.Close()

	// Prepare a WebSocket connection request
	u := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Make a request to the test server
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Error making a request to the test server: %v", err)
	}
	defer conn.Close()

	// Send a test message to the server
	testMessage := "Hello, server!"
	err = conn.WriteMessage(websocket.TextMessage, []byte(testMessage))
	if err != nil {
		t.Fatalf("Error sending message to the server: %v", err)
	}

	// Read the response from the server
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading message from the server: %v", err)
	}

	// Check if the response matches the sent message
	if string(message) != testMessage {
		t.Errorf("Expected message %q, got %q", testMessage, string(message))
	}
}

func TestHandleWebSocketDisconnection(t *testing.T) {
	var clientDisconnected int32 = 0
	atomic.StoreInt32(&clientDisconnected, 0)

	// Create an HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		handleWebSocketDisconnection(conn, &clientDisconnected)
	}))
	defer ts.Close()

	// Prepare a WebSocket connection request
	u := "ws" + strings.TrimPrefix(ts.URL, "http")

	// Make a request to the test server
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("Error making a request to the test server: %v", err)
	}

	// Close the connection to simulate a disconnection
	conn.Close()

	// Wait for the server to handle the disconnection
	time.Sleep(time.Millisecond * 100)

	// Check if the server has handled the disconnection
	if atomic.LoadInt32(&clientDisconnected) != 1 {
		t.Error("Server did not handle the client disconnection")
	}
}
