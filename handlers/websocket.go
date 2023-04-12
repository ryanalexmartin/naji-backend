package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
}

func handleWebSocketMessage(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read message from the WebSocket connection
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Echo the message back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}

func handleWebSocketDisconnection(conn *websocket.Conn, clientDisconnected *int32) {
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			atomic.StoreInt32(clientDisconnected, 1)
			break
		}
	}
}
