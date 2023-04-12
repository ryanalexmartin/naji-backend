package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	ID   string
	Conn MessageWriter
	//	Conn *websocket.Conn
	Room string
}

func NewUser(conn *websocket.Conn) *User {
	return &User{
		ID:   uuid.New().String(),
		Conn: conn,
	}
}

type MessageWriter interface {
	WriteMessage(messageType int, data []byte) error
}
